package handshake

import "github.com/flynn/noise"

var cipherSuite = noise.NewCipherSuite(noise.DH25519, noise.CipherChaChaPoly, noise.HashSHA256)

func getNoiseConfig(handshakeRole role, publicKey []byte, privateKey []byte) noise.Config {
	return noise.Config{
		Pattern:     noise.HandshakeNK,
		Initiator:   handshakeRole == Initiator,
		CipherSuite: cipherSuite,
		StaticKeypair: func() noise.DHKey {
			if handshakeRole == Responder {
				return noise.DHKey{Private: privateKey, Public: publicKey}
			} else {
				return noise.DHKey{}
			}
		}(),
		PeerStatic: func() []byte {
			if handshakeRole == Initiator {
				return publicKey
			} else {
				return nil
			}
		}(),
	}
}
