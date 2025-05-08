package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/GrahamZiervogel/pokedex/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *Config, args ...string) error
}

type Config struct {
	PokeapiClient            *pokeapi.Client
	NextLocationAreasURL     *string
	PreviousLocationAreasURL *string
	Pokedex                  map[string]pokeapi.Pokemon
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	httpClientTimeout := 5 * time.Second
	cacheReapInterval := 5 * time.Minute

	pokeClient := pokeapi.NewClient(httpClientTimeout, cacheReapInterval)

	cfg := &Config{
		PokeapiClient: pokeClient,
		Pokedex:       make(map[string]pokeapi.Pokemon),
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
		args := []string{}
		if len(cleanedWords) > 1 {
			args = cleanedWords[1:]
		}

		command, exists := getCommands()[commandName]
		if exists {
			err := command.callback(cfg, args...)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
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
		"explore": {
			name:        "explore <location_area_name>",
			description: "Lists Pokémon in a given location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch <pokemon_name>",
			description: "Attempt to catch a Pokémon and add it to your Pokedex",
			callback:    commandCatch,
		},
	}
}
