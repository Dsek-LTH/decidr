package handshake

import "github.com/flynn/noise"

var cipherSuite = noise.NewCipherSuite(noise.DH25519, noise.CipherChaChaPoly, noise.HashSHA256)

func getNoiseConfig(handshakeRole role) noise.Config {
	return noise.Config{
		Pattern:     noise.HandshakeNN,
		Initiator:   handshakeRole == Initiator,
		CipherSuite: cipherSuite,
	}
}
