package client

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/therewardstore/httpmatter"
)

// mockHTTPClient is a mock implementation of HTTPClient for testing.
type mockHTTPClient struct {
	doFunc func(url string) (*http.Response, error)
}

func (m *mockHTTPClient) Get(url string) (*http.Response, error) {
	return m.doFunc(url)
}

func TestFetchData(t *testing.T) {
	// Initialize httpmatter to ensure HTTP calls are mocked
	httpmatter.Init(&httpmatter.Config{})
	_ = httpmatter.NewHTTP(t)

	tests := []struct {
		name        string
		setupMock   func() HTTPClient
		wantErr     bool
		errContains string
		validate    func(*testing.T, []byte)
	}{
		{
			name: "successful fetch with valid JSON response",
			setupMock: func() HTTPClient {
				return &mockHTTPClient{
					doFunc: func(url string) (*http.Response, error) {
						if url != Endpoint {
							t.Errorf("expected URL %s, got %s", Endpoint, url)
						}
						// JSON Placeholder style response
						jsonResponse := `{
							"id": 1,
							"name": "Bangalore",
							"country": "India",
							"population": 8443675,
							"state": "Karnataka"
						}`
						resp := &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(strings.NewReader(jsonResponse)),
							Header:     make(http.Header),
						}
						resp.Header.Set("Content-Type", "application/json")
						return resp, nil
					},
				}
			},
			wantErr: false,
			validate: func(t *testing.T, body []byte) {
				if len(body) == 0 {
					t.Error("expected non-empty response body")
				}
				bodyStr := string(body)
				if !strings.Contains(bodyStr, "Bangalore") {
					t.Errorf("expected response to contain 'Bangalore', got: %s", bodyStr)
				}
				if !strings.Contains(bodyStr, "India") {
					t.Errorf("expected response to contain 'India', got: %s", bodyStr)
				}
			},
		},
		{
			name: "HTTP 404 Not Found",
			setupMock: func() HTTPClient {
				return &mockHTTPClient{
					doFunc: func(url string) (*http.Response, error) {
						resp := &http.Response{
							StatusCode: http.StatusNotFound,
							Body:       io.NopCloser(strings.NewReader(`{"error": "Not Found"}`)),
							Header:     make(http.Header),
						}
						resp.Header.Set("Content-Type", "application/json")
						return resp, nil
					},
				}
			},
			wantErr:     true,
			errContains: "unexpected status code: 404",
		},
		{
			name: "HTTP 500 Internal Server Error",
			setupMock: func() HTTPClient {
				return &mockHTTPClient{
					doFunc: func(url string) (*http.Response, error) {
						resp := &http.Response{
							StatusCode: http.StatusInternalServerError,
							Body:       io.NopCloser(strings.NewReader(`{"error": "Internal Server Error"}`)),
							Header:     make(http.Header),
						}
						resp.Header.Set("Content-Type", "application/json")
						return resp, nil
					},
				}
			},
			wantErr:     true,
			errContains: "unexpected status code: 500",
		},
		{
			name: "HTTP 503 Service Unavailable",
			setupMock: func() HTTPClient {
				return &mockHTTPClient{
					doFunc: func(url string) (*http.Response, error) {
						resp := &http.Response{
							StatusCode: http.StatusServiceUnavailable,
							Body:       io.NopCloser(strings.NewReader(`{"error": "Service Unavailable"}`)),
							Header:     make(http.Header),
						}
						resp.Header.Set("Content-Type", "application/json")
						return resp, nil
					},
				}
			},
			wantErr:     true,
			errContains: "unexpected status code: 503",
		},
		{
			name: "network error - connection refused",
			setupMock: func() HTTPClient {
				return &mockHTTPClient{
					doFunc: func(url string) (*http.Response, error) {
						return nil, &http.ProtocolError{
							ErrorString: "connection refused",
						}
					},
				}
			},
			wantErr:     true,
			errContains: "failed to fetch data",
		},
		{
			name: "empty response body",
			setupMock: func() HTTPClient {
				return &mockHTTPClient{
					doFunc: func(url string) (*http.Response, error) {
						resp := &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(strings.NewReader("")),
							Header:     make(http.Header),
						}
						resp.Header.Set("Content-Type", "application/json")
						return resp, nil
					},
				}
			},
			wantErr: false,
			validate: func(t *testing.T, body []byte) {
				if len(body) != 0 {
					t.Errorf("expected empty body, got: %s", string(body))
				}
			},
		},
		{
			name: "malformed JSON response",
			setupMock: func() HTTPClient {
				return &mockHTTPClient{
					doFunc: func(url string) (*http.Response, error) {
						resp := &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(strings.NewReader(`{"name": "Bangalore", invalid json`)),
							Header:     make(http.Header),
						}
						resp.Header.Set("Content-Type", "application/json")
						return resp, nil
					},
				}
			},
			wantErr: false, // We just return the body, parsing is not our responsibility
			validate: func(t *testing.T, body []byte) {
				if len(body) == 0 {
					t.Error("expected non-empty response body")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := tt.setupMock()
			got, err := FetchData(mockClient)

			if (err != nil) != tt.wantErr {
				t.Errorf("FetchData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.errContains != "" && err != nil {
					if !strings.Contains(err.Error(), tt.errContains) {
						t.Errorf("FetchData() error = %v, should contain %v", err, tt.errContains)
					}
				}
				if got != nil {
					t.Errorf("FetchData() = %v, want nil on error", got)
				}
			} else {
				if err != nil {
					t.Errorf("FetchData() unexpected error = %v", err)
					return
				}
				if tt.validate != nil {
					tt.validate(t, got)
				}
			}
		})
	}
}

func TestDefaultClient_Get(t *testing.T) {
	client := NewDefaultClient()

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid URL",
			url:     "https://example.com",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.Get(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultClient.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if resp != nil {
				resp.Body.Close()
			}
		})
	}
}

func TestNewDefaultClient(t *testing.T) {
	client := NewDefaultClient()
	if client == nil {
		t.Error("NewDefaultClient() returned nil")
	}
	if client.client == nil {
		t.Error("NewDefaultClient() client is nil")
	}
}
