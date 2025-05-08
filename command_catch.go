package main

import (
	"errors"
	"fmt"
	"math/rand"
)

func commandCatch(cfg *Config, args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide exactly one Pokémon name to catch")
	}
	pokemonName := args[0]

	if _, caught := cfg.Pokedex[pokemonName]; caught {
		fmt.Printf("%s is already in your Pokedex!\n", pokemonName)
		return nil
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemonData, err := cfg.PokeapiClient.GetPokemonDetails(pokemonName)
	if err != nil {
		return err
	}

	const maxRollValue = 500
	const minCatchScore = 50
	const maxCatchScore = maxRollValue - 50

	catchScore := maxRollValue - pokemonData.BaseExperience

	if catchScore < minCatchScore {
		catchScore = minCatchScore
	}
	if catchScore > maxCatchScore {
		catchScore = maxCatchScore
	}

	roll := rand.Intn(maxRollValue)

	if roll < catchScore {
		fmt.Printf("%s was caught!\n", pokemonData.Name)
		cfg.Pokedex[pokemonData.Name] = pokemonData
		fmt.Printf("%s added to Pokedex.\n", pokemonData.Name)
	} else {
		fmt.Printf("%s escaped!\n", pokemonData.Name)
	}

	return nil
}
