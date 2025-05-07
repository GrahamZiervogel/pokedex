package main

import (
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	trimmedText := strings.TrimSpace(text)
	lowercasedText := strings.ToLower(trimmedText)

	words := strings.Fields(lowercasedText)

	return words
}

func main() {
	text := "  Hello   World  "
	cleanedWords := cleanInput(text)
	fmt.Println("Cleaned words:", cleanedWords)
}
