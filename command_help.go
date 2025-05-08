package main

import (
	"errors"
	"fmt"
	"sort"
)

func commandHelp(cfg *Config, args ...string) error {
	if len(args) > 0 {
		return errors.New("help command does not take any arguments")
	}

	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	commands := getCommands()
	commandNames := make([]string, 0, len(commands))

	for name := range commands {
		commandNames = append(commandNames, name)
	}

	sort.Strings(commandNames)

	for _, name := range commandNames {
		cmd := commands[name]
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	fmt.Println()
	return nil
}
