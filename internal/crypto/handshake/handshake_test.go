package handshake

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Dsek-LTH/decidr/internal/crypto"
)

func TestHandshakeVerificationWordsMatch(t *testing.T) {
	client, server := newInMemoryPeers()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	clientIdentity, adminIdentity := getIdentityPair(t)
	clientCh := runHandshakeAsync(ctx, client, clientIdentity)
	serverCh := runHandshakeAsync(ctx, server, adminIdentity)

	var clientRes, serverRes handshakeResult

	select {
	case clientRes = <-clientCh:
	case <-ctx.Done():
		t.Fatal("client handshake timed out")
	}

	select {
	case serverRes = <-serverCh:
	case <-ctx.Done():
		t.Fatal("server handshake timed out")
	}

	if clientRes.err != nil {
		t.Fatalf("client handshake failed: %v", clientRes.err)
	}
	if serverRes.err != nil {
		t.Fatalf("server handshake failed: %v", serverRes.err)
	}

	if len(clientRes.hash) == 0 {
		t.Fatal("client handshake hash is empty")
	}
	if len(serverRes.hash) == 0 {
		t.Fatal("server handshake hash is empty")
	}

	clientWords := crypto.GetVerificationWords(clientRes.hash, 6)
	serverWords := crypto.GetVerificationWords(serverRes.hash, 6)

	if strings.Join(clientWords, "-") != strings.Join(serverWords, "-") {
		t.Fatalf(
			"verification codes differ:\nclient: %v\nserver: %v",
			clientWords,
			serverWords,
		)
	}
}

func TestHandshakeContextCancellation(t *testing.T) {
	client, server := newInMemoryPeers()

	ctx, cancel := context.WithCancel(context.Background())

	clientIdentity, adminIdentity := getIdentityPair(t)
	clientCh := runHandshakeAsync(ctx, client, clientIdentity)
	serverCh := runHandshakeAsync(ctx, server, adminIdentity)

	// Cancel shortly after starting
	cancel()

	clientRes := <-clientCh
	serverRes := <-serverCh

	if clientRes.err == nil && serverRes.err == nil {
		t.Fatal("expected handshake to fail due to context cancellation")
	}
}

func TestHandshakeSendFailure(t *testing.T) {
	client, server := newInMemoryPeers()

	client.sendFunc = func(context.Context, []byte) error {
		return errors.New("send failed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	clientIdentity, adminIdentity := getIdentityPair(t)
	clientCh := runHandshakeAsync(ctx, client, clientIdentity)
	serverCh := runHandshakeAsync(ctx, server, adminIdentity)

	clientRes := <-clientCh
	if clientRes.err == nil {
		t.Fatal("expected client handshake to fail")
	}

	// Cancel to unblock responder
	cancel()

	select {
	case <-serverCh:
		// Responder should also fail due to send failure on client side
	case <-time.After(time.Second):
		t.Fatal("responder did not exit after context cancellation")
	}
}

func TestHandshakeHashMatchesBetweenPeers(t *testing.T) {
	client, server := newInMemoryPeers()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	clientIdentity, adminIdentity := getIdentityPair(t)
	clientCh := runHandshakeAsync(ctx, client, clientIdentity)
	serverCh := runHandshakeAsync(ctx, server, adminIdentity)

	clientRes := <-clientCh
	serverRes := <-serverCh

	if !bytes.Equal(clientRes.hash, serverRes.hash) {
		t.Fatal("handshake hashes do not match")
	}
}

func TestMultipleHandshakesVerificationWordsMatch(t *testing.T) {
	const n = 5

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	results := make(chan error, n)

	clientIdentity, adminIdentity := getIdentityPair(t)
	var verificationWords []string
	for range n {
		go func() {
			client, server := newInMemoryPeers()

			clientCh := runHandshakeAsync(ctx, client, clientIdentity)
			serverCh := runHandshakeAsync(ctx, server, adminIdentity)

			clientRes := <-clientCh
			serverRes := <-serverCh

			if clientRes.err != nil {
				results <- clientRes.err
				return
			}
			if serverRes.err != nil {
				results <- serverRes.err
				return
			}
			if !bytes.Equal(clientRes.hash, serverRes.hash) {
				results <- errors.New("hash mismatch")
				return
			}
			verificationWords = append(verificationWords, strings.Join(
				crypto.GetVerificationWords(clientRes.hash, 6),
				"-",
			))

			results <- nil
		}()
	}

	for range n {
		if err := <-results; err != nil {
			t.Fatal(err)
		}
	}

	first := verificationWords[0]
	for _, words := range verificationWords[1:] {
		if words != first {
			t.Fatal("verification words do not match across handshakes")
		}
	}
}
