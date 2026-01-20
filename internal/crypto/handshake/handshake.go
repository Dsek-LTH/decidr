// Package handshake contains functions to perform Noise Protocol Framework handshakes.
package handshake

import (
	"context"

	"github.com/flynn/noise"
)

// Perform executes a Noise handshake using the specified role (Initiator or Responder).
//
// Returns the send and receive cipher states for encrypting/decrypting subsequent messages,
// the handshake hash for verification, and/or any error encountered during the process.
//
// sendCipherState is used for outbound traffic after the handshake.
//
// receiveCipherState is used for inbound traffic after the handshake.
func Perform(
	ctx context.Context,
	send func([]byte) error,
	receive func() ([]byte, error),
	identity handshakeIdentity,
) (
	sendCipherState *noise.CipherState,
	receiveCipherState *noise.CipherState,
	handshakeState *noise.HandshakeState,
	err error,
) {
	peer := NewFuncPeer(send, receive)
	handshakeState, err = noise.NewHandshakeState(identity.getNoiseConfig())
	if err != nil {
		return nil, nil, nil, err
	}

	var cipherState1, cipherState2 *noise.CipherState
	for _, step := range stepsFor(identity.getRole()) {
		cipherState1, cipherState2, err = step.apply(ctx, handshakeState, peer)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	sendCipherState, receiveCipherState = cipherStatesFor(
		identity.getRole(),
		cipherState1,
		cipherState2,
	)

	// sendCipherState is used for outbound traffic after the handshake.
	// receiveCipherState is used for inbound traffic after the handshake.
	return sendCipherState, receiveCipherState, handshakeState, nil
}
