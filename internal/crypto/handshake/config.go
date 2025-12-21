package handshake

import "github.com/flynn/noise"

var cipherSuite = noise.NewCipherSuite(noise.DH25519, noise.CipherChaChaPoly, noise.HashSHA256)

func GetNoiseConfig(handshakeRole Roles) noise.Config {
	return noise.Config{
		Pattern:     noise.HandshakeNN,
		Initiator:   handshakeRole == Initiator,
		CipherSuite: cipherSuite,
	}
}
