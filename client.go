package client

import (
	"fmt"
	"io"
	"net/http"
)

// Test CI trigger - small change to verify pipeline

const (
	// Endpoint is the API endpoint to fetch data from
	// Using JSON Placeholder API to get all posts
	Endpoint = "https://jsonplaceholder.typicode.com/posts"
)

// HTTPClient defines the interface for HTTP operations.
// This allows dependency injection for testing.
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// DefaultClient wraps the standard http.Client.
type DefaultClient struct {
	client *http.Client
}

// NewDefaultClient creates a new default HTTP client.
func NewDefaultClient() *DefaultClient {
	return &DefaultClient{
		client: &http.Client{},
	}
}

// Get performs an HTTP GET request.
func (c *DefaultClient) Get(url string) (*http.Response, error) {
	return c.client.Get(url)
}

// FetchData fetches data from the endpoint and processes the response.
// It accepts an HTTPClient interface for dependency injection, making it easily testable.
// The response is expected to be in JSON format (similar to JSON Placeholder API).
func FetchData(client HTTPClient) ([]byte, error) {
	// Make HTTP GET request
	resp, err := client.Get(Endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read and process response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}
