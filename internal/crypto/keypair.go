package crypto

import (
	"crypto/ed25519"
	"crypto/rand"

	"golang.org/x/crypto/curve25519"
)

func GenerateStaticKeypair() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	var privateKey [32]byte
	_, err := rand.Read(privateKey[:])
	if err != nil {
		return nil, nil, err
	}
	publicKey, err := curve25519.X25519(privateKey[:], curve25519.Basepoint)
	if err != nil {
		return nil, nil, err
	}

	return publicKey, privateKey[:], nil
}
