package pokeapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	timeout := 5 * time.Second
	reapInterval := 10 * time.Minute
	client := NewClient(timeout, reapInterval)

	if client.httpClient.Timeout != timeout {
		t.Errorf("expected http client timeout %v, got %v", timeout, client.httpClient.Timeout)
	}
	if client.cache == nil {
		t.Errorf("expected cache to be initialized, but it was nil")
	}
}

func TestListLocationAreas_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/location-area" {
			t.Errorf("Expected to request '/location-area', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"count":1,"next":"https://next.url","previous":null,"results":[{"name":"test-location","url":"https://location.url"}]}`)
	}))
	defer server.Close()

	client := NewClient(5*time.Second, 5*time.Minute)
	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	var pageURL *string

	resp, err := client.ListLocationAreas(pageURL)
	if err != nil {
		t.Fatalf("ListLocationAreas failed: %v", err)
	}

	if resp.Count != 1 {
		t.Errorf("expected count 1, got %d", resp.Count)
	}
	if resp.Results[0].Name != "test-location" {
		t.Errorf("expected location name 'test-location', got '%s'", resp.Results[0].Name)
	}
	if resp.Next == nil || *resp.Next != "https://next.url" {
		t.Errorf("expected next URL 'https://next.url', got '%v'", resp.Next)
	}
}

func TestListLocationAreas_CacheHit(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"count":1,"next":null,"previous":null,"results":[{"name":"cached-location","url":"https://cached.url"}]}`)
	}))
	defer server.Close()

	client := NewClient(5*time.Second, 100*time.Millisecond)
	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	_, err := client.ListLocationAreas(nil)
	if err != nil {
		t.Fatalf("First call to ListLocationAreas failed: %v", err)
	}
	if requestCount != 1 {
		t.Errorf("expected 1 server request after first call, got %d", requestCount)
	}

	resp, err := client.ListLocationAreas(nil)
	if err != nil {
		t.Fatalf("Second call to ListLocationAreas failed: %v", err)
	}
	if requestCount != 1 {
		t.Errorf("expected 1 server request total (cache hit), got %d", requestCount)
	}
	if resp.Results[0].Name != "cached-location" {
		t.Errorf("expected location name 'cached-location' from cache, got '%s'", resp.Results[0].Name)
	}
}

func TestListLocationAreas_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(5*time.Second, 5*time.Minute)
	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	_, err := client.ListLocationAreas(nil)
	if err == nil {
		t.Fatal("ListLocationAreas did not return an error for API failure")
	}
}

func TestListLocationAreas_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"count":1, "results": [{"name":"test"`)
	}))
	defer server.Close()

	client := NewClient(5*time.Second, 5*time.Minute)
	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	_, err := client.ListLocationAreas(nil)
	if err == nil {
		t.Fatal("ListLocationAreas did not return an error for malformed JSON")
	}
}

func TestListLocationAreas_PaginationURL(t *testing.T) {
	expectedPathOnServer := "/specific-page-url"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != expectedPathOnServer {
			t.Errorf("Expected server to receive request for path '%s', got: '%s'", expectedPathOnServer, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"count":1,"next":null,"previous":null,"results":[{"name":"page-location","url":"someurl"}]}`)
	}))
	defer server.Close()

	client := NewClient(5*time.Second, 5*time.Minute)

	pageURLToTest := server.URL + expectedPathOnServer

	_, err := client.ListLocationAreas(&pageURLToTest)
	if err != nil {
		t.Fatalf("ListLocationAreas with pageURL failed: %v", err)
	}
}
