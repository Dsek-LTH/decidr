package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/Dsek-LTH/decidr/internal/crypto/handshake"
	"github.com/gorilla/websocket"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		runAdmin()
		wg.Done()
	}()
	go func() {
		runClient()
		wg.Done()
	}()

	wg.Wait()
}

func runClient() {
	ctx := context.Background()

	conn, _, err := websocket.DefaultDialer.Dial(
		"ws://localhost:8080/ws/client?id=client-1&admin=admin-1",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	transportPeer := handshake.NewFuncPeer(
		func(b []byte) error {
			return conn.WriteMessage(websocket.BinaryMessage, b)
		},
		func() ([]byte, error) {
			_, msg, err := conn.ReadMessage()
			return msg, err
		},
	)

	send := func(b []byte) error {
		return transportPeer.Send(ctx, b)
	}

	receive := func() ([]byte, error) {
		return transportPeer.Receive(ctx)
	}

	msg, _ := transportPeer.Receive(ctx)
	clientEndpoint := handshake.GetClientEndpoint(msg)
	fmt.Println("[client] client endpoint identity:", clientEndpoint.Identity)

	sendCS, recvCS, _, err := handshake.Perform(
		ctx,
		send,
		receive,
		clientEndpoint.Identity,
	)
	if err != nil {
		log.Fatal("[client] handshake failed:", err)
	} else {
		fmt.Println("[client] handshake succeeded")
	}

	securePeer := handshake.NewFuncPeer(
		// Send: encrypt → transport
		func(plaintext []byte) error {
			ciphertext, err := sendCS.Encrypt(nil, nil, plaintext)
			if err != nil {
				return err
			}
			return transportPeer.Send(ctx, ciphertext)
		},
		// Receive: transport → decrypt
		func() ([]byte, error) {
			ciphertext, err := transportPeer.Receive(ctx)
			if err != nil {
				return nil, err
			}
			return recvCS.Decrypt(nil, nil, ciphertext)
		},
	)

	// Send message
	_ = securePeer.Send(ctx, []byte("hello admin"))
	fmt.Println("[client] message sent: hello admin")

	// Receive response
	msg, _ = securePeer.Receive(ctx)
	fmt.Println("[client] message received:", string(msg))
}

func runAdmin() {
	ctx := context.Background()

	conn, _, err := websocket.DefaultDialer.Dial(
		"ws://localhost:8080/ws/admin?id=admin-1",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	transportPeer := handshake.NewFuncPeer(
		func(b []byte) error {
			return conn.WriteMessage(websocket.BinaryMessage, b)
		},
		func() ([]byte, error) {
			_, msg, err := conn.ReadMessage()
			return msg, err
		},
	)

	send := func(b []byte) error {
		return transportPeer.Send(ctx, append([]byte("client-1\n"), b...))
	}

	receive := func() ([]byte, error) {
		return transportPeer.Receive(ctx)
	}

	clientEndpoint, adminEndpoint, _ := handshake.NewAdminEndpoint()
	fmt.Println("[admin] client endpoint identity:", clientEndpoint.Identity)

	if err := transportPeer.Send(
		ctx,
		append([]byte("client-1\n"), clientEndpoint.Identity.GetPublicKey()...),
	); err != nil {
		log.Fatal("[admin] failed to send public key to client:", err)
	}

	sendCS, recvCS, _, err := handshake.Perform(
		ctx,
		send,
		receive,
		adminEndpoint.Identity,
	)
	if err != nil {
		log.Fatal("[admin] handshake failed:", err)
	} else {
		fmt.Println("[admin] handshake succeeded")
	}

	securePeer := handshake.NewFuncPeer(
		// Send: encrypt → transport
		func(plaintext []byte) error {
			ciphertext, err := sendCS.Encrypt(nil, nil, plaintext)
			if err != nil {
				return err
			}
			return transportPeer.Send(ctx, append([]byte("client-1\n"), ciphertext...))
		},
		// Receive: transport → decrypt
		func() ([]byte, error) {
			ciphertext, err := transportPeer.Receive(ctx)
			if err != nil {
				return nil, err
			}
			return recvCS.Decrypt(nil, nil, ciphertext)
		},
	)

	// Receive message
	msg, _ := securePeer.Receive(ctx)
	fmt.Println("[admin] message received:", string(msg))

	// Send response
	_ = securePeer.Send(ctx, []byte("hello client"))
	fmt.Println("[admin] message sent: hello client")
}
