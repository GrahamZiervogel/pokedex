package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/GrahamZiervogel/pokedex/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *Config) error
}

type Config struct {
	Next     *string
	Previous *string
	Cache    *pokecache.Cache
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	cacheInterval := 5 * time.Minute
	cfg := &Config{
		Cache: pokecache.NewCache(cacheInterval),
	}

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "Error reading input:", err)
			}
			fmt.Println("\nExiting Pokedex REPL.")
			break
		}

		userInput := scanner.Text()
		cleanedWords := cleanInput(userInput)

		if len(cleanedWords) == 0 {
			continue
		}

		commandName := cleanedWords[0]
		command, exists := getCommands()[commandName]
		if exists {
			err := command.callback(cfg)
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}

func cleanInput(text string) []string {
	trimmedText := strings.TrimSpace(text)
	lowercasedText := strings.ToLower(trimmedText)

	words := strings.Fields(lowercasedText)

	return words
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display the next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 location areas",
			callback:    commandMapb,
		},
	}
}
