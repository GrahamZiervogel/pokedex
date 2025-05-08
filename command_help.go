package main

import (
	"errors"
	"fmt"
)

func commandHelp(cfg *Config, args ...string) error {
	if len(args) > 0 {
		return errors.New("help command does not take any arguments")
	}
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()
	return nil
}
