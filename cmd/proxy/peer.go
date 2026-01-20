package main

import (
	"github.com/Dsek-LTH/decidr/internal/crypto/handshake"
	"github.com/gorilla/websocket"
)

func newWebSocketPeer(conn *websocket.Conn) handshake.Peer {
	return handshake.NewFuncPeer(
		func(b []byte) error {
			return conn.WriteMessage(websocket.BinaryMessage, b)
		},
		func() ([]byte, error) {
			_, msg, err := conn.ReadMessage()
			return msg, err
		},
	)
}
