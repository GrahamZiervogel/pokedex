package pokeapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetLocationAreaDetails_Success(t *testing.T) {
	areaName := "test-area"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := fmt.Sprintf("/location-area/%s", areaName)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected to request path '%s', got '%s'", expectedPath, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{
			"name": "test-area",
			"pokemon_encounters": [
				{"pokemon": {"name": "pidgey"}},
				{"pokemon": {"name": "rattata"}}
			]
		}`)
	}))
	defer server.Close()

	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	client := NewClient(5*time.Second, 5*time.Minute)
	details, err := client.GetLocationAreaDetails(areaName)

	if err != nil {
		t.Fatalf("GetLocationAreaDetails failed: %v", err)
	}
	if details.Name != areaName {
		t.Errorf("Expected area name '%s', got '%s'", areaName, details.Name)
	}
	if len(details.PokemonEncounters) != 2 {
		t.Fatalf("Expected 2 pokemon encounters, got %d", len(details.PokemonEncounters))
	}
	if details.PokemonEncounters[0].Pokemon.Name != "pidgey" {
		t.Errorf("Expected first pokemon 'pidgey', got '%s'", details.PokemonEncounters[0].Pokemon.Name)
	}
	if details.PokemonEncounters[1].Pokemon.Name != "rattata" {
		t.Errorf("Expected second pokemon 'rattata', got '%s'", details.PokemonEncounters[1].Pokemon.Name)
	}
}

func TestGetLocationAreaDetails_CacheHit(t *testing.T) {
	areaName := "cache-test-area"
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		expectedPath := fmt.Sprintf("/location-area/%s", areaName)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected to request path '%s', got '%s'", expectedPath, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"name": "%s", "pokemon_encounters": [{"pokemon": {"name": "bulbasaur"}}]}`, areaName)
	}))
	defer server.Close()

	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	client := NewClient(5*time.Second, 5*time.Minute)

	_, err := client.GetLocationAreaDetails(areaName)
	if err != nil {
		t.Fatalf("First call to GetLocationAreaDetails failed: %v", err)
	}
	if requestCount != 1 {
		t.Errorf("Expected 1 server request after first call, got %d", requestCount)
	}

	details, err := client.GetLocationAreaDetails(areaName)
	if err != nil {
		t.Fatalf("Second call to GetLocationAreaDetails failed: %v", err)
	}
	if requestCount != 1 {
		t.Errorf("Expected 1 server request total (cache hit), got %d", requestCount)
	}
	if details.PokemonEncounters[0].Pokemon.Name != "bulbasaur" {
		t.Errorf("Expected 'bulbasaur' from cache, got '%s'", details.PokemonEncounters[0].Pokemon.Name)
	}
}

func TestGetLocationAreaDetails_NotFound(t *testing.T) {
	areaName := "nonexistent-area"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	client := NewClient(5*time.Second, 5*time.Minute)
	_, err := client.GetLocationAreaDetails(areaName)

	if err == nil {
		t.Fatal("Expected an error for 404 Not Found, but got nil")
	}
	expectedErrorMsg := fmt.Sprintf("location area '%s' not found", areaName)
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedErrorMsg, err.Error())
	}
}

func TestGetLocationAreaDetails_APIError(t *testing.T) {
	areaName := "error-area"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	client := NewClient(5*time.Second, 5*time.Minute)
	_, err := client.GetLocationAreaDetails(areaName)

	if err == nil {
		t.Fatal("Expected an error for API server error, but got nil")
	}
	if !strings.Contains(err.Error(), "API request to") || !strings.Contains(err.Error(), "failed with status code: 500") {
		t.Errorf("Error message '%s' did not match expected format for 500 error", err.Error())
	}
}

func TestGetLocationAreaDetails_MalformedJSON(t *testing.T) {
	areaName := "malformed-json-area"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"name": "test-area", "pokemon_encounters": [`)
	}))
	defer server.Close()

	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	client := NewClient(5*time.Second, 5*time.Minute)
	_, err := client.GetLocationAreaDetails(areaName)

	if err == nil {
		t.Fatal("Expected an error for malformed JSON, but got nil")
	}
	if !strings.Contains(err.Error(), "error unmarshalling JSON") {
		t.Errorf("Error message '%s' did not match expected format for unmarshalling error", err.Error())
	}
}

func TestGetLocationAreaDetails_EmptyArgument(t *testing.T) {
	client := NewClient(5*time.Second, 5*time.Minute)
	_, err := client.GetLocationAreaDetails("")
	if err == nil {
		t.Fatal("Expected an error for empty area name, but got nil")
	}
	expectedErrorMsg := "location area name or ID cannot be empty"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMsg, err.Error())
	}
}
