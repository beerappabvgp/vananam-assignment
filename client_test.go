package client

import (
	"fmt"
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
						// JSON Placeholder API response format - array of posts
						jsonResponse := `[
							{
								"userId": 1,
								"id": 1,
								"title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
								"body": "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto"
							},
							{
								"userId": 1,
								"id": 2,
								"title": "qui est esse",
								"body": "est rerum tempore vitae\nsequi sint nihil reprehenderit dolor beatae ea dolores neque"
							}
						]`
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
				if !strings.Contains(bodyStr, "userId") {
					t.Errorf("expected response to contain 'userId', got: %s", bodyStr)
				}
				if !strings.Contains(bodyStr, "title") {
					t.Errorf("expected response to contain 'title', got: %s", bodyStr)
				}
				// Verify it's an array (starts with [)
				if !strings.HasPrefix(strings.TrimSpace(bodyStr), "[") {
					t.Errorf("expected array response, got: %s", bodyStr[:50])
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
							Body:       io.NopCloser(strings.NewReader(`[{"userId": 1, invalid json`)),
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
		{
			name: "HTTP 400 Bad Request",
			setupMock: func() HTTPClient {
				return &mockHTTPClient{
					doFunc: func(url string) (*http.Response, error) {
						resp := &http.Response{
							StatusCode: http.StatusBadRequest,
							Body:       io.NopCloser(strings.NewReader(`{"error": "Bad Request"}`)),
							Header:     make(http.Header),
						}
						resp.Header.Set("Content-Type", "application/json")
						return resp, nil
					},
				}
			},
			wantErr:     true,
			errContains: "unexpected status code: 400",
		},
		{
			name: "HTTP 401 Unauthorized",
			setupMock: func() HTTPClient {
				return &mockHTTPClient{
					doFunc: func(url string) (*http.Response, error) {
						resp := &http.Response{
							StatusCode: http.StatusUnauthorized,
							Body:       io.NopCloser(strings.NewReader(`{"error": "Unauthorized"}`)),
							Header:     make(http.Header),
						}
						resp.Header.Set("Content-Type", "application/json")
						return resp, nil
					},
				}
			},
			wantErr:     true,
			errContains: "unexpected status code: 401",
		},
		{
			name: "HTTP 403 Forbidden",
			setupMock: func() HTTPClient {
				return &mockHTTPClient{
					doFunc: func(url string) (*http.Response, error) {
						resp := &http.Response{
							StatusCode: http.StatusForbidden,
							Body:       io.NopCloser(strings.NewReader(`{"error": "Forbidden"}`)),
							Header:     make(http.Header),
						}
						resp.Header.Set("Content-Type", "application/json")
						return resp, nil
					},
				}
			},
			wantErr:     true,
			errContains: "unexpected status code: 403",
		},
		{
			name: "HTTP 502 Bad Gateway",
			setupMock: func() HTTPClient {
				return &mockHTTPClient{
					doFunc: func(url string) (*http.Response, error) {
						resp := &http.Response{
							StatusCode: http.StatusBadGateway,
							Body:       io.NopCloser(strings.NewReader(`{"error": "Bad Gateway"}`)),
							Header:     make(http.Header),
						}
						resp.Header.Set("Content-Type", "application/json")
						return resp, nil
					},
				}
			},
			wantErr:     true,
			errContains: "unexpected status code: 502",
		},
		{
			name: "network error - timeout",
			setupMock: func() HTTPClient {
				return &mockHTTPClient{
					doFunc: func(url string) (*http.Response, error) {
						return nil, &http.ProtocolError{
							ErrorString: "timeout",
						}
					},
				}
			},
			wantErr:     true,
			errContains: "failed to fetch data",
		},
		{
			name: "network error - DNS resolution failed",
			setupMock: func() HTTPClient {
				return &mockHTTPClient{
					doFunc: func(url string) (*http.Response, error) {
						return nil, &http.ProtocolError{
							ErrorString: "no such host",
						}
					},
				}
			},
			wantErr:     true,
			errContains: "failed to fetch data",
		},
		{
			name: "response body read error",
			setupMock: func() HTTPClient {
				return &mockHTTPClient{
					doFunc: func(url string) (*http.Response, error) {
						// Create a response with a body that will error on read
						resp := &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(&errorReader{}),
							Header:     make(http.Header),
						}
						resp.Header.Set("Content-Type", "application/json")
						return resp, nil
					},
				}
			},
			wantErr:     true,
			errContains: "failed to read response body",
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

// errorReader is a reader that always returns an error
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("read error")
}
