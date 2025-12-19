package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
)

func GenerateIdentityKeypair() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)

	return publicKey, privateKey, err
}
