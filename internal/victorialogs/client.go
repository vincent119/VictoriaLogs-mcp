package victorialogs

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/vincent119/victorialogs-mcp/internal/util"
)

// Client VictoriaLogs HTTP client
type Client struct {
	httpClient *util.HTTPClient
	baseURL    string
	maxResults int
}

// ClientOption client option
type ClientOption func(*Client)

// WithMaxResults sets max results
func WithMaxResults(max int) ClientOption {
	return func(c *Client) {
		c.maxResults = max
	}
}

// NewClient creates new VictoriaLogs client
func NewClient(baseURL string, auth util.AuthConfig, timeout time.Duration, opts ...ClientOption) *Client {
	c := &Client{
		baseURL:    baseURL,
		maxResults: 5000,
	}

	c.httpClient = util.NewHTTPClient(
		util.WithBaseURL(baseURL),
		util.WithAuth(auth),
		util.WithTimeout(timeout),
	)

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Close closes client
func (c *Client) Close() {
	if c.httpClient != nil {
		c.httpClient.Close()
	}
}

// doRequest executes HTTP request
func (c *Client) doRequest(ctx context.Context, method, path string, query url.Values) ([]byte, error) {
	fullPath := path
	if len(query) > 0 {
		fullPath = path + "?" + query.Encode()
	}

	resp, err := c.httpClient.Do(ctx, method, fullPath, nil)
	if err != nil {
		return nil, &APIError{
			StatusCode: 0,
			Message:    fmt.Sprintf("HTTP request failed: %v", err),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	return body, nil
}

// Health health check
func (c *Client) Health(ctx context.Context) (*HealthResponse, error) {
	resp, err := c.httpClient.Get(ctx, "/health")
	if err != nil {
		return nil, &APIError{
			StatusCode: 0,
			Message:    fmt.Sprintf("health check failed: %v", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &HealthResponse{Status: "unhealthy"}, nil
	}

	return &HealthResponse{Status: "healthy"}, nil
}

// GetMaxResults gets max results setting
func (c *Client) GetMaxResults() int {
	return c.maxResults
}
