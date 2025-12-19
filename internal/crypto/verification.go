// Package crypto provides cryptographic utilities for the Decidr application.
package crypto

import (
	"github.com/Dsek-LTH/decidr/pkg/data/wordlists"
)

func GetVerificationWords(hash []byte, wordCount int) []string {
	bitLength := len(hash) * 8
	bits := make([]bool, bitLength)

	for i := range bitLength {
		bits[i] = (hash[i/8] & (1 << (7 - (i % 8)))) != 0
	}

	words := []string{}
	for i := 0; i+11 <= bitLength && len(words) < wordCount; i += 11 {
		val := 0
		for j := range 11 {
			val <<= 1
			if bits[i+j] {
				val |= 1
			}
		}
		words = append(words, wordlists.GetBip39Wordlist(wordlists.English)[val])
	}

	return words
}
