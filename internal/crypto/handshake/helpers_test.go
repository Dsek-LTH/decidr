package handshake

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"testing"

	"golang.org/x/crypto/curve25519"
)

type inMemoryPeer struct {
	sendFunc    func(ctx context.Context, b []byte) error
	receiveFunc func(ctx context.Context) ([]byte, error)
	sendCh      chan []byte
	receiveCh   chan []byte
}

func (p inMemoryPeer) Send(ctx context.Context, b []byte) error {
	if p.sendFunc != nil {
		return p.sendFunc(ctx, b)
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case p.sendCh <- b:
		return nil
	}
}

func (p inMemoryPeer) Receive(ctx context.Context) ([]byte, error) {
	if p.receiveFunc != nil {
		return p.receiveFunc(ctx)
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case msg, ok := <-p.receiveCh:
		if !ok {
			return nil, context.Canceled
		}
		return msg, nil
	}
}

func newInMemoryPeers() (client, server inMemoryPeer) {
	c2s := make(chan []byte, 1)
	s2c := make(chan []byte, 1)

	clientPeer := inMemoryPeer{sendCh: c2s, receiveCh: s2c}
	serverPeer := inMemoryPeer{sendCh: s2c, receiveCh: c2s}

	return clientPeer, serverPeer
}

type handshakeResult struct {
	hash []byte
	err  error
}

func runHandshakeAsync(ctx context.Context, role role, peer peer, publicKey []byte, privateKey []byte) <-chan handshakeResult {
	ch := make(chan handshakeResult, 1)

	go func() {
		send := func(b []byte) error { return peer.Send(ctx, b) }
		receive := func() ([]byte, error) { return peer.Receive(ctx) }

		if _, _, handshakeState, err := Perform(ctx, role, send, receive, publicKey, privateKey); err != nil {
			ch <- handshakeResult{err: fmt.Errorf("handshake failed: %w", err)}
		} else {
			sum := func() [32]byte {
				if role == Initiator {
					return sha256.Sum256(handshakeState.PeerStatic())
				} else {
					return sha256.Sum256(publicKey)
				}
			}()
			ch <- handshakeResult{hash: sum[:], err: err}
		}
	}()

	return ch
}

func getKeypair(t *testing.T) (ed25519.PublicKey, ed25519.PrivateKey) {
	t.Helper()

	var privateKey [32]byte
	_, err := rand.Read(privateKey[:])
	if err != nil {
		t.Fatalf("failed to read random bytes for X25519 private key: %v", err)
	}
	publicKey, err := curve25519.X25519(privateKey[:], curve25519.Basepoint)
	if err != nil {
		t.Fatalf("failed to derive X25519 public key: %v", err)
	}

	return publicKey, privateKey[:]
}
