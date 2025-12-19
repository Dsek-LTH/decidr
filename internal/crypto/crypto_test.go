package crypto

import (
	"crypto/ed25519"
	"fmt"
	"strings"
)

func ExampleGetVerificationWords() {
	// Actual usage would be:
	// publicKey, _, _ := GenerateIdentityKeypair()

	// but for this example, we use a fixed value to get a deterministic output.
	publicKey := ed25519.PublicKey("0")
	hash := GetVerificationHash(publicKey, publicKey)
	words := GetVerificationWords(hash, 6)
	code := strings.Join(words, "-")

	fmt.Printf("Verification Code: %s\n", code)
	// output: Verification Code: vanish-old-tonight-execute-saddle-thank
}
