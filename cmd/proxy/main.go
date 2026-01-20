package main

import (
	"log"
	"net/http"

	"github.com/Dsek-LTH/decidr/internal/crypto/handshake"
	"github.com/gorilla/websocket"
)

var (
	router   = handshake.NewRouter()
	upgrader = websocket.Upgrader{}
)

func main() {
	http.HandleFunc("/ws/client", clientHandler)
	http.HandleFunc("/ws/admin", adminHandler)

	log.Println("Proxy server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
