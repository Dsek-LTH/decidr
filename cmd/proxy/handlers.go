package main

import (
	"log"
	"net/http"
)

func clientHandler(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("id")
	adminID := r.URL.Query().Get("admin")
	log.Println("New client connection:", clientID, "->", adminID)

	if clientID == "" || adminID == "" {
		http.Error(w, "missing id or admin", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	peer := newWebSocketPeer(conn)
	router.RegisterClient(clientID, peer)
	defer router.RemoveClient(clientID)

	ctx := r.Context()

	for {
		msg, err := peer.Receive(ctx)
		if err != nil {
			return
		}
		log.Println("client msg:", string(msg))

		// Forward client → admin
		if err := router.RouteToAdmin(ctx, adminID, msg); err != nil {
			log.Println("route error:", err)
		}
	}
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	adminID := r.URL.Query().Get("id")
	if adminID == "" {
		http.Error(w, "missing admin id", http.StatusBadRequest)
		return
	}
	log.Println("New admin connection:", adminID)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	peer := newWebSocketPeer(conn)
	router.RegisterAdmin(adminID, peer)
	defer router.RemoveAdmin(adminID)

	ctx := r.Context()

	for {
		msg, err := peer.Receive(ctx)
		if err != nil {
			return
		}

		// Example protocol:
		// "<client-id>\n<payload>"
		clientID, payload := split(msg)
		if clientID == "" {
			continue
		}
		log.Println("admin msg:", string(payload))

		if err := router.RouteToClient(ctx, clientID, payload); err != nil {
			log.Println("route error:", err)
		}
	}
}

func split(msg []byte) (string, []byte) {
	for i, b := range msg {
		if b == '\n' {
			return string(msg[:i]), msg[i+1:]
		}
	}
	return "", nil
}
