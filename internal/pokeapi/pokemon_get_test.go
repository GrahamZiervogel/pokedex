package pokeapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetPokemonDetails_Success(t *testing.T) {
	pokemonName := "pikachu"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := fmt.Sprintf("/pokemon/%s", pokemonName)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected to request path '%s', got '%s'", expectedPath, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{
			"id": 25,
			"name": "pikachu",
			"base_experience": 112,
			"height": 4,
			"weight": 60,
			"stats": [
				{"base_stat": 35, "stat": {"name": "hp"}},
				{"base_stat": 55, "stat": {"name": "attack"}}
			],
			"types": [
				{"slot": 1, "type": {"name": "electric"}}
			]
		}`)
	}))
	defer server.Close()

	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	client := NewClient(5*time.Second, 5*time.Minute)
	details, err := client.GetPokemonDetails(pokemonName)

	if err != nil {
		t.Fatalf("GetPokemonDetails failed: %v", err)
	}

	if details.Name != pokemonName {
		t.Errorf("Expected pokemon name '%s', got '%s'", pokemonName, details.Name)
	}
	if details.ID != 25 {
		t.Errorf("Expected pokemon ID 25, got %d", details.ID)
	}
	if details.BaseExperience != 112 {
		t.Errorf("Expected base experience 112, got %d", details.BaseExperience)
	}
	if details.Height != 4 {
		t.Errorf("Expected height 4, got %d", details.Height)
	}
	if details.Weight != 60 {
		t.Errorf("Expected weight 60, got %d", details.Weight)
	}
	if len(details.Types) != 1 || details.Types[0].Type.Name != "electric" {
		t.Errorf("Expected type 'electric', got %v", details.Types)
	}
	if len(details.Stats) != 2 || details.Stats[0].Stat.Name != "hp" || details.Stats[0].BaseStat != 35 {
		t.Errorf("Expected first stat 'hp' with base 35, got %v", details.Stats)
	}
}

func TestGetPokemonDetails_CacheHit(t *testing.T) {
	pokemonName := "charmander"
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		expectedPath := fmt.Sprintf("/pokemon/%s", pokemonName)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected to request path '%s', got '%s'", expectedPath, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"id": 4, "name": "%s", "base_experience": 62}`, pokemonName)
	}))
	defer server.Close()

	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	client := NewClient(5*time.Second, 5*time.Minute)

	_, err := client.GetPokemonDetails(pokemonName)
	if err != nil {
		t.Fatalf("First call to GetPokemonDetails failed: %v", err)
	}
	if requestCount != 1 {
		t.Errorf("Expected 1 server request after first call, got %d", requestCount)
	}

	details, err := client.GetPokemonDetails(pokemonName)
	if err != nil {
		t.Fatalf("Second call to GetPokemonDetails failed: %v", err)
	}
	if requestCount != 1 {
		t.Errorf("Expected 1 server request total (cache hit), got %d", requestCount)
	}
	if details.Name != pokemonName || details.BaseExperience != 62 {
		t.Errorf("Expected '%s' with base_experience 62 from cache, got '%s' with %d", pokemonName, details.Name, details.BaseExperience)
	}
}

func TestGetPokemonDetails_NotFound(t *testing.T) {
	pokemonName := "nonexistentpokemon"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == fmt.Sprintf("/pokemon/%s", pokemonName) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Unexpected request to %s", r.URL.Path)
		}
	}))
	defer server.Close()

	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	client := NewClient(5*time.Second, 5*time.Minute)
	_, err := client.GetPokemonDetails(pokemonName)

	if err == nil {
		t.Fatal("Expected an error for 404 Not Found, but got nil")
	}
	expectedErrorMsg := fmt.Sprintf("pokemon '%s' not found", pokemonName)
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedErrorMsg, err.Error())
	}
}

func TestGetPokemonDetails_APIError(t *testing.T) {
	pokemonName := "errorpokemon"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	client := NewClient(5*time.Second, 5*time.Minute)
	_, err := client.GetPokemonDetails(pokemonName)

	if err == nil {
		t.Fatal("Expected an error for API server error, but got nil")
	}

	actualError := err.Error()

	if !strings.Contains(actualError, "API request to") {
		t.Errorf("Error message '%s' did not contain 'API request to'", actualError)
	}

	expectedStatusCodePart := fmt.Sprintf("failed with status code %d:", http.StatusInternalServerError)
	if !strings.Contains(actualError, expectedStatusCodePart) {
		t.Errorf("Error message '%s' did not contain '%s'", actualError, expectedStatusCodePart)
	}
}

func TestGetPokemonDetails_MalformedJSON(t *testing.T) {
	pokemonName := "badjsonpokemon"
	malformedJSON := `{"id": 1, "name": "badjsonpokemon", "base_experience": 100, "types": [` // Intentionally broken

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, malformedJSON)
	}))
	defer server.Close()

	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	client := NewClient(5*time.Second, 5*time.Minute)
	_, err := client.GetPokemonDetails(pokemonName)

	if err == nil {
		t.Fatal("Expected an error for malformed JSON, but got nil")
	}
	if !strings.Contains(err.Error(), "error unmarshalling JSON") {
		t.Errorf("Error message '%s' did not match expected format for unmarshalling error", err.Error())
	}
	if !strings.Contains(err.Error(), malformedJSON) {
		t.Errorf("Error message '%s' did not include the malformed JSON response body", err.Error())
	}
}

func TestGetPokemonDetails_EmptyArgument(t *testing.T) {
	client := NewClient(5*time.Second, 5*time.Minute)
	_, err := client.GetPokemonDetails("")
	if err == nil {
		t.Fatal("Expected an error for empty Pokémon name, but got nil")
	}
	expectedErrorMsg := "pokemon name cannot be empty"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMsg, err.Error())
	}
}
