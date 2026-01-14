package client

import (
	"fmt"
	"io"
	"net/http"
)

const (
	Endpoint = "https://jsonplaceholder.typicode.com/posts"
)

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type DefaultClient struct {
	client *http.Client
}

func NewDefaultClient() *DefaultClient {
	return &DefaultClient{
		client: &http.Client{},
	}
}

func (c *DefaultClient) Get(url string) (*http.Response, error) {
	return c.client.Get(url)
}

func FetchData(client HTTPClient) ([]byte, error) {
	resp, err := client.Get(Endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}
