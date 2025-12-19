// Package noise provides functionality related to the Noise Protocol Framework.
package noise

import "github.com/flynn/noise"

type HandshakeRoles int

const (
	Initiator HandshakeRoles = iota
	Responder
)

var cipherSuite = noise.NewCipherSuite(noise.DH25519, noise.CipherChaChaPoly, noise.HashSHA256)

func GetNoiseConfig(handshakeRole HandshakeRoles) noise.Config {
	return noise.Config{
		Pattern:     noise.HandshakeNN,
		Initiator:   handshakeRole == Initiator,
		CipherSuite: cipherSuite,
	}
}
