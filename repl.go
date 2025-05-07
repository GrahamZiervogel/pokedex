package main

import (
	"strings"
)

func cleanInput(text string) []string {
	trimmedText := strings.TrimSpace(text)
	lowercasedText := strings.ToLower(trimmedText)

	words := strings.Fields(lowercasedText)

	return words
}
