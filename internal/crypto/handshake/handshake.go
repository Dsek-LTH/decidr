// Package handshake contains functions to perform Noise Protocol Framework handshakes.
package handshake

import (
	"context"
	"errors"

	"github.com/flynn/noise"
)

// Perform executes a Noise handshake using the specified role (Initiator or Responder).
//
// Returns the send and receive cipher states for encrypting/decrypting subsequent messages,
// the handshake hash for verification, and/or any error encountered during the process.
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
	if !role.valid() {
		return nil, nil, nil, errors.New("invalid handshake role")
	}

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

	// sendCipherState is used for outbound traffic after the handshake.
	// receiveCipherState is used for inbound traffic after the handshake.
	return sendCipherState, receiveCipherState, handshakeState.ChannelBinding(), nil
}
