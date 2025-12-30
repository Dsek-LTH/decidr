package handshake

import "github.com/flynn/noise"

type role int

const (
	initiator role = iota
	responder
)

func cipherStatesFor(role role, cipherState1, cipherState2 *noise.CipherState) (sendCipherState, receiveCipherState *noise.CipherState) {
	if role == initiator {
		// For Initiator, cipherState1 is Send, cipherState2 is Receive
		return cipherState1, cipherState2
	} else {
		// For Responder, cipherState1 is Receive, cipherState2 is Send
		return cipherState2, cipherState1
	}
}
