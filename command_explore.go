package main

import (
	"errors"
	"fmt"
)

func commandExplore(cfg *Config, args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide exactly one location area name to explore")
	}
	locationAreaName := args[0]

	fmt.Printf("Exploring %s...\n", locationAreaName)

	areaDetails, err := cfg.PokeapiClient.GetLocationAreaDetails(locationAreaName)
	if err != nil {
		return fmt.Errorf("could not get details for %s: %w", locationAreaName, err)
	}

	if len(areaDetails.PokemonEncounters) == 0 {
		fmt.Printf("No Pokémon found in %s.\n", locationAreaName)
		return nil
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range areaDetails.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}
