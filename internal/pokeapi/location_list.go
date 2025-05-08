package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) ListLocationAreas(pageURL *string) (LocationAreaResponse, error) {
	url := BaseURL + "/location-area"
	if pageURL != nil && *pageURL != "" {
		url = *pageURL
	}

	var responseBody []byte
	var locationResponse LocationAreaResponse

	cachedData, found := c.cache.Get(url)
	if found {
		responseBody = cachedData
	} else {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return LocationAreaResponse{}, fmt.Errorf("error creating HTTP request: %w", err)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return LocationAreaResponse{}, fmt.Errorf("error making HTTP request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode > 299 {
			return LocationAreaResponse{}, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, url)
		}

		responseBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return LocationAreaResponse{}, fmt.Errorf("error reading response body: %w", err)
		}

		c.cache.Add(url, responseBody)
	}

	err := json.Unmarshal(responseBody, &locationResponse)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("error unmarshalling JSON for %s: %w", url, err)
	}

	return locationResponse, nil
}
