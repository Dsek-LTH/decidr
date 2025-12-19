// Package wordlists provides access to embedded wordlists for various languages.
package wordlists

import (
	"embed"
	"io/fs"
	"strings"
)

//go:embed *.txt
var WordlistsFS embed.FS

type WordlistLanguages int

const (
	English WordlistLanguages = iota
)

func GetBip39Wordlist(language WordlistLanguages) []string {
	wordlistBytes, _ := fs.ReadFile(WordlistsFS, "bip39_english.txt")
	wordlistStr := string(wordlistBytes)

	return strings.Split(wordlistStr, "\n")
}
