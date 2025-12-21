package handshake

import "github.com/flynn/noise"

type step interface {
	apply(*noise.HandshakeState, peer) (*noise.CipherState, *noise.CipherState, error)
}

type stepSend struct{}

func (stepSend) apply(
	handshakeState *noise.HandshakeState,
	peer peer,
) (*noise.CipherState, *noise.CipherState, error) {
	message, sendCipherState, receiveCipherState, err := handshakeState.WriteMessage(nil, nil)
	if err != nil {
		return nil, nil, err
	}

	if err := peer.Send(message); err != nil {
		return nil, nil, err
	}

	return sendCipherState, receiveCipherState, nil
}

type stepReceive struct{}

func (stepReceive) apply(
	handshakeState *noise.HandshakeState,
	peer peer,
) (*noise.CipherState, *noise.CipherState, error) {
	message, err := peer.Receive()
	if err != nil {
		return nil, nil, err
	}

	_, sendCipherState, receiveCipherState, err := handshakeState.ReadMessage(nil, message)
	return sendCipherState, receiveCipherState, err
}
