package handshake

import (
	"context"
	"fmt"
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

func runHandshakeAsync(ctx context.Context, role role, peer peer) <-chan handshakeResult {
	ch := make(chan handshakeResult, 1)

	go func() {
		send := func(b []byte) error { return peer.Send(ctx, b) }
		receive := func() ([]byte, error) { return peer.Receive(ctx) }

		if _, _, hash, err := Perform(ctx, role, send, receive); err != nil {
			ch <- handshakeResult{err: fmt.Errorf("handshake failed: %w", err)}
		} else {
			ch <- handshakeResult{hash: hash, err: err}
		}
	}()

	return ch
}
