package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://pokeapi.co/api/v2/location-area"

type LocationAreaResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func commandMap(cfg *Config) error {
	url := baseURL

	if cfg.Next != nil && *cfg.Next != "" {
		url = *cfg.Next
	}

	return DisplayLocationAreas(cfg, url)
}

func commandMapb(cfg *Config) error {
	if cfg.Previous == nil || *cfg.Previous == "" {
		fmt.Println("You are at the first page of locations.")
		return nil
	}

	url := *cfg.Previous

	return DisplayLocationAreas(cfg, url)
}

func DisplayLocationAreas(cfg *Config, url string) error {
	resp, err := http.Get(url)

	if err != nil {
		return fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	var locationResponse LocationAreaResponse
	err = json.Unmarshal(body, &locationResponse)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	cfg.Next = locationResponse.Next
	cfg.Previous = locationResponse.Previous

	for _, area := range locationResponse.Results {
		fmt.Printf("%s\n", area.Name)
	}
	if len(locationResponse.Results) == 0 {
		fmt.Println("No location areas found.")
	}

	return nil
}
