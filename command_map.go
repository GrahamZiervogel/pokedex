package main

import (
	"errors"
	"fmt"
)

func commandMap(cfg *Config, args ...string) error {
	if len(args) > 0 {
		return errors.New("map command does not take any arguments")
	}

	fmt.Println("Fetching next location areas...")
	locationResponse, err := cfg.PokeapiClient.ListLocationAreas(cfg.NextLocationAreasURL)
	if err != nil {
		return fmt.Errorf("could not get location areas: %w", err)
	}

	cfg.NextLocationAreasURL = locationResponse.Next
	cfg.PreviousLocationAreasURL = locationResponse.Previous

	if len(locationResponse.Results) == 0 {
		fmt.Println("No more location areas found.")
		return nil
	}

	fmt.Println("Location Areas:")
	for _, area := range locationResponse.Results {
		fmt.Printf("- %s\n", area.Name)
	}

	return nil
}

func commandMapb(cfg *Config, args ...string) error {
	if len(args) > 0 {
		return errors.New("mapb command does not take any arguments")
	}

	if cfg.PreviousLocationAreasURL == nil || *cfg.PreviousLocationAreasURL == "" {
		fmt.Println("You are at the first page of locations, cannot go back.")
		return nil
	}

	fmt.Println("Fetching previous location areas...")
	locationResponse, err := cfg.PokeapiClient.ListLocationAreas(cfg.PreviousLocationAreasURL)
	if err != nil {
		return fmt.Errorf("could not get previous location areas: %w", err)
	}

	cfg.NextLocationAreasURL = locationResponse.Next
	cfg.PreviousLocationAreasURL = locationResponse.Previous

	if len(locationResponse.Results) == 0 {
		fmt.Println("No location areas found on the previous page.")
		return nil
	}

	fmt.Println("Location Areas (Previous):")
	for _, area := range locationResponse.Results {
		fmt.Printf("- %s\n", area.Name)
	}

	return nil
}
