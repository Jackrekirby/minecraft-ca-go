package core

import (
	"strings"
)

func ToSnakeCase(input string) string {
	// Split the string by spaces
	words := strings.Fields(input)

	// Convert each word to lowercase
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}

	// Join the words with underscores
	return strings.Join(words, "_")
}
