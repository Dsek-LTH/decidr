package handshake

import (
	"context"
	"strings"
	"testing"

	"github.com/Dsek-LTH/decidr/internal/crypto"
)

func TestHandshakeVerificationWordsMatch(t *testing.T) {
	// In-memory "network"
	clientToServerChannel := make(chan []byte, 1)
	serverToClientChannel := make(chan []byte, 1)

	// Mock send/receive
	clientSend := func(b []byte) error {
		clientToServerChannel <- b
		return nil
	}
	clientReceive := func() ([]byte, error) {
		return <-serverToClientChannel, nil
	}

	serverSend := func(b []byte) error {
		serverToClientChannel <- b
		return nil
	}
	serverReceive := func() ([]byte, error) {
		return <-clientToServerChannel, nil
	}

	type result struct {
		hash []byte
		err  error
	}

	clientResultChannel := make(chan result, 1)
	serverResultChannel := make(chan result, 1)

	// Run client and server concurrently
	go func() {
		_, _, hash, err := Perform(context.Background(), Initiator, clientSend, clientReceive)
		clientResultChannel <- result{hash: hash, err: err}
	}()

	go func() {
		_, _, hash, err := Perform(context.Background(), Responder, serverSend, serverReceive)
		serverResultChannel <- result{hash: hash, err: err}
	}()

	clientResult := <-clientResultChannel
	serverResult := <-serverResultChannel

	if clientResult.err != nil {
		t.Fatalf("client handshake failed: %v", clientResult.err)
	}
	if serverResult.err != nil {
		t.Fatalf("server handshake failed: %v", serverResult.err)
	}

	clientWords := crypto.GetVerificationWords(clientResult.hash, 6)
	serverWords := crypto.GetVerificationWords(serverResult.hash, 6)

	clientCode := strings.Join(clientWords, "-")
	serverCode := strings.Join(serverWords, "-")

	if clientCode != serverCode {
		t.Fatalf(
			"verification codes differ:\nclient: %s\nserver: %s",
			clientCode,
			serverCode,
		)
	}
}

func TestHandshakeContextCancellation(t *testing.T) {
	// In-memory "network"
	clientToServerChannel := make(chan []byte, 1)
	serverToClientChannel := make(chan []byte, 1)

	// Mock send/receive
	clientSend := func(b []byte) error {
		clientToServerChannel <- b
		return nil
	}
	clientReceive := func() ([]byte, error) {
		return <-serverToClientChannel, nil
	}

	serverSend := func(b []byte) error {
		serverToClientChannel <- b
		return nil
	}
	serverReceive := func() ([]byte, error) {
		return <-clientToServerChannel, nil
	}

	// Create a context that is already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, _, _, err := Perform(ctx, Initiator, clientSend, clientReceive)
	if err == nil {
		t.Fatal("expected client handshake to fail due to context cancellation, but it succeeded")
	}

	_, _, _, err = Perform(ctx, Responder, serverSend, serverReceive)
	if err == nil {
		t.Fatal("expected server handshake to fail due to context cancellation, but it succeeded")
	}
}
