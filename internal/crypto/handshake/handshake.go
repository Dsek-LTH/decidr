// Package handshake contains functions to perform Noise Protocol Framework handshakes.
package handshake

import (
	"context"

	"github.com/flynn/noise"
)

func Perform(
	ctx context.Context,
	role role,
	send func([]byte) error,
	receive func() ([]byte, error),
) (
	sendCipherState *noise.CipherState,
	receiveCipherState *noise.CipherState,
	handshakeHash []byte,
	err error,
) {
	peer := newFuncPeer(send, receive)

	handshakeState, err := noise.NewHandshakeState(
		getNoiseConfig(role),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	for _, step := range stepsFor(role) {
		sendCipherState, receiveCipherState, err = step.apply(ctx, handshakeState, peer)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return sendCipherState, receiveCipherState, handshakeState.ChannelBinding(), nil
}
