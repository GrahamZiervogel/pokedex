package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetLocationAreaDetails(locationAreaNameOrID string) (LocationAreaDetailsResponse, error) {
	if locationAreaNameOrID == "" {
		return LocationAreaDetailsResponse{}, fmt.Errorf("location area name or ID cannot be empty")
	}

	url := fmt.Sprintf("%s/location-area/%s", BaseURL, locationAreaNameOrID)

	var responseBody []byte
	var areaDetails LocationAreaDetailsResponse

	cachedData, found := c.cache.Get(url)
	if found {
		responseBody = cachedData
	} else {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return LocationAreaDetailsResponse{}, fmt.Errorf("error creating HTTP request for %s: %w", url, err)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return LocationAreaDetailsResponse{}, fmt.Errorf("error making HTTP request to %s: %w", url, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			return LocationAreaDetailsResponse{}, fmt.Errorf("location area '%s' not found", locationAreaNameOrID)
		}
		if resp.StatusCode > 299 {
			return LocationAreaDetailsResponse{}, fmt.Errorf("API request to %s failed with status code: %d", url, resp.StatusCode)
		}

		responseBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return LocationAreaDetailsResponse{}, fmt.Errorf("error reading response body from %s: %w", url, err)
		}

		c.cache.Add(url, responseBody)
	}

	err := json.Unmarshal(responseBody, &areaDetails)
	if err != nil {
		return LocationAreaDetailsResponse{}, fmt.Errorf("error unmarshalling JSON for %s: %w", url, err)
	}

	return areaDetails, nil
}
