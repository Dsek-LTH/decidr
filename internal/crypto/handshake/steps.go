package handshake

import (
	"context"
	"fmt"

	"github.com/flynn/noise"
)

type step interface {
	apply(
		context.Context,
		*noise.HandshakeState,
		Peer,
	) (*noise.CipherState, *noise.CipherState, error)
}

type stepSend struct{}

func (stepSend) apply(
	ctx context.Context,
	handshakeState *noise.HandshakeState,
	peer Peer,
) (*noise.CipherState, *noise.CipherState, error) {
	message, cipherState1, cipherState2, err := handshakeState.WriteMessage(nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("handshake write: %w", err)
	}

	if err := peer.Send(ctx, message); err != nil {
		return nil, nil, fmt.Errorf("handshake send: %w", err)
	}

	return cipherState1, cipherState2, nil
}

type stepReceive struct{}

func (stepReceive) apply(
	ctx context.Context,
	handshakeState *noise.HandshakeState,
	peer Peer,
) (*noise.CipherState, *noise.CipherState, error) {
	message, err := peer.Receive(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("handshake receive: %w", err)
	}

	_, cipherState1, cipherState2, err := handshakeState.ReadMessage(nil, message)
	if err != nil {
		err = fmt.Errorf("handshake read: %w", err)
	}
	return cipherState1, cipherState2, err
}
