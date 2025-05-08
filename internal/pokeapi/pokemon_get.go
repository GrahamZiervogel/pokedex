package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetPokemonDetails(pokemonName string) (Pokemon, error) {
	if pokemonName == "" {
		return Pokemon{}, fmt.Errorf("pokemon name cannot be empty")
	}

	url := fmt.Sprintf("%s/pokemon/%s", BaseURL, pokemonName)

	var responseBody []byte
	var pokemonDetails Pokemon

	cachedData, found := c.cache.Get(url)
	if found {
		responseBody = cachedData
	} else {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return Pokemon{}, fmt.Errorf("error creating HTTP request for %s: %w", url, err)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return Pokemon{}, fmt.Errorf("error making HTTP request to %s: %w", url, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			return Pokemon{}, fmt.Errorf("pokemon '%s' not found", pokemonName)
		}
		if resp.StatusCode > 299 {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return Pokemon{}, fmt.Errorf("API request to %s failed with status code %d: %s", url, resp.StatusCode, string(bodyBytes))
		}

		responseBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return Pokemon{}, fmt.Errorf("error reading response body from %s: %w", url, err)
		}

		c.cache.Add(url, responseBody)
	}

	err := json.Unmarshal(responseBody, &pokemonDetails)
	if err != nil {
		return Pokemon{}, fmt.Errorf("error unmarshalling JSON for %s: %w (Response: %s)", url, err, string(responseBody))
	}

	return pokemonDetails, nil
}
