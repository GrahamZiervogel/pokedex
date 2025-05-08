package main

import (
	"errors"
	"fmt"
	"sort"
)

func commandPokedex(cfg *Config, args ...string) error {
	if len(args) > 0 {
		return errors.New("pokedex command does not take any arguments")
	}

	fmt.Println("Your Pokedex:")

	if len(cfg.Pokedex) == 0 {
		fmt.Println(" (is empty)")
		return nil
	}

	pokemonNames := make([]string, 0, len(cfg.Pokedex))
	for name := range cfg.Pokedex {
		pokemonNames = append(pokemonNames, name)
	}
	sort.Strings(pokemonNames)

	for _, name := range pokemonNames {
		fmt.Printf(" - %s\n", name)
	}

	return nil
}
