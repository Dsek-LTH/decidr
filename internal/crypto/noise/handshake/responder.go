package handshake

import (
	internalNoise "github.com/Dsek-LTH/decidr/internal/crypto/noise"
	"github.com/flynn/noise"
)

func ResponderHandshake(send func([]byte) error, receive func() ([]byte, error)) (
	sendCipherState *noise.CipherState,
	receiveCipherState *noise.CipherState,
	handshakeHash []byte,
	err error,
) {
	handshakeState, err := noise.NewHandshakeState(internalNoise.GetNoiseConfig(internalNoise.Responder))
	if err != nil {
		return
	}

	// ← Message 1 (client → server)
	msg1, err := receive()
	if err != nil {
		return
	}

	_, _, _, err = handshakeState.ReadMessage(nil, msg1)
	if err != nil {
		return
	}

	// → Message 2 (server → client)
	msg2, sendCipherState, receiveCipherState, err := handshakeState.WriteMessage(nil, nil)
	if err != nil {
		return
	}

	err = send(msg2)
	if err != nil {
		return
	}

	handshakeHash = handshakeState.ChannelBinding()

	return
}
