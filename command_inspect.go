package main

import (
	"errors"
	"fmt"
)

func commandInspect(cfg *Config, args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide exactly one Pokémon name to inspect")
	}
	pokemonName := args[0]

	pokemon, caught := cfg.Pokedex[pokemonName]
	if !caught {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, statEntry := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", statEntry.Stat.Name, statEntry.BaseStat)
	}

	fmt.Println("Types:")
	for _, typeEntry := range pokemon.Types {
		fmt.Printf("  - %s\n", typeEntry.Type.Name)
	}

	return nil
}
