// Package handshake contains functions to perform Noise Protocol Framework handshakes.
package handshake

import (
	"github.com/flynn/noise"
)

type Roles int

const (
	Initiator Roles = iota
	Responder
)

func Perform(
	role Roles,
	send func([]byte) error,
	receive func() ([]byte, error),
) (
	sendCipherState *noise.CipherState,
	receiveCipherState *noise.CipherState,
	handshakeHash []byte,
	err error,
) {
	handshakeState, err := noise.NewHandshakeState(GetNoiseConfig(role))
	if err != nil {
		return
	}

	var steps []stepFunction
	if role == Initiator {
		steps = []stepFunction{
			stepSend(send),
			stepReceive(receive),
		}
	} else {
		steps = []stepFunction{
			stepReceive(receive),
			stepSend(send),
		}
	}
	for _, step := range steps {
		sendCipherState, receiveCipherState, err = step(handshakeState)
		if err != nil {
			return
		}
	}

	handshakeHash = handshakeState.ChannelBinding()
	return
}
