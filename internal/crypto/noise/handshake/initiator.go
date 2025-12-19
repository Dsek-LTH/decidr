// Package handshake contains functions to perform Noise Protocol Framework handshakes.
package handshake

import (
	internalNoise "github.com/Dsek-LTH/decidr/internal/crypto/noise"
	"github.com/flynn/noise"
)

func InitiatorHandshake(send func([]byte) error, receive func() ([]byte, error)) (
	sendCipherState *noise.CipherState,
	receiveCipherState *noise.CipherState,
	handshakeHash []byte,
	err error,
) {
	handshakeState, err := noise.NewHandshakeState(internalNoise.GetNoiseConfig(internalNoise.Initiator))
	if err != nil {
		return
	}

	// -> Message 1 (client -> server)
	msg1, _, _, err := handshakeState.WriteMessage(nil, nil)
	if err != nil {
		return
	}

	err = send(msg1)
	if err != nil {
		return
	}

	// <- Message 2 (server -> client)
	msg2, err := receive()
	if err != nil {
		return
	}

	_, sendCipherState, receiveCipherState, err = handshakeState.ReadMessage(nil, msg2)
	if err != nil {
		return
	}

	handshakeHash = handshakeState.ChannelBinding()

	return
}
