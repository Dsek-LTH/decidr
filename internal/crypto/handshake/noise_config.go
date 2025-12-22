package handshake

import (
	"github.com/flynn/noise"
)

var cipherSuite = noise.NewCipherSuite(noise.DH25519, noise.CipherChaChaPoly, noise.HashSHA256)

func (c clientIdentity) getNoiseConfig() noise.Config {
	return noise.Config{
		Pattern:     noise.HandshakeNK,
		Initiator:   true,
		CipherSuite: cipherSuite,
		PeerStatic:  c.AdminPublicKey,
	}
}

func (a adminIdentity) getNoiseConfig() noise.Config {
	return noise.Config{
		Pattern:     noise.HandshakeNK,
		Initiator:   false,
		CipherSuite: cipherSuite,
		StaticKeypair: noise.DHKey{
			Public:  a.PublicKey,
			Private: a.PrivateKey,
		},
	}
}
