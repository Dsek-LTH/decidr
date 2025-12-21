package handshake

import "github.com/flynn/noise"

type stepFunction func(*noise.HandshakeState) (*noise.CipherState, *noise.CipherState, error)

func stepSend(send func([]byte) error) stepFunction {
	return func(handshakeState *noise.HandshakeState) (*noise.CipherState, *noise.CipherState, error) {
		message, sendCipherState, receiveCipherState, err := handshakeState.WriteMessage(nil, nil)
		if err != nil {
			return nil, nil, err
		}
		if err := send(message); err != nil {
			return nil, nil, err
		}

		return sendCipherState, receiveCipherState, nil
	}
}

func stepReceive(receive func() ([]byte, error)) stepFunction {
	return func(handshakeState *noise.HandshakeState) (*noise.CipherState, *noise.CipherState, error) {
		message, err := receive()
		if err != nil {
			return nil, nil, err
		}

		_, sendCipherState, receiveCipherState, err := handshakeState.ReadMessage(nil, message)

		return sendCipherState, receiveCipherState, err
	}
}
