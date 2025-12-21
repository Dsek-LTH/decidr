package handshake

import (
	"context"

	"github.com/flynn/noise"
)

type step interface {
	apply(
		context.Context,
		*noise.HandshakeState,
		peer,
	) (*noise.CipherState, *noise.CipherState, error)
}

type stepSend struct{}

func (stepSend) apply(
	ctx context.Context,
	handshakeState *noise.HandshakeState,
	peer peer,
) (*noise.CipherState, *noise.CipherState, error) {
	message, sendCipherState, receiveCipherState, err := handshakeState.WriteMessage(nil, nil)
	if err != nil {
		return nil, nil, err
	}

	if err := peer.Send(ctx, message); err != nil {
		return nil, nil, err
	}

	return sendCipherState, receiveCipherState, nil
}

type stepReceive struct{}

func (stepReceive) apply(
	ctx context.Context,
	handshakeState *noise.HandshakeState,
	peer peer,
) (*noise.CipherState, *noise.CipherState, error) {
	message, err := peer.Receive(ctx)
	if err != nil {
		return nil, nil, err
	}

	_, sendCipherState, receiveCipherState, err := handshakeState.ReadMessage(nil, message)
	return sendCipherState, receiveCipherState, err
}
