// Package util 提供通用工具函式
package util

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

// HTTPClient HTTP 客戶端封裝
type HTTPClient struct {
	client  *http.Client
	baseURL string
	auth    AuthConfig
}

// AuthConfig 認證設定
type AuthConfig struct {
	Type     string // none | basic | bearer
	Username string
	Password string
	Token    string
}

// HTTPClientOption HTTP 客戶端選項
type HTTPClientOption func(*HTTPClient)

// WithBaseURL 設定 Base URL
func WithBaseURL(url string) HTTPClientOption {
	return func(c *HTTPClient) {
		c.baseURL = url
	}
}

// WithAuth 設定認證
func WithAuth(auth AuthConfig) HTTPClientOption {
	return func(c *HTTPClient) {
		c.auth = auth
	}
}

// WithTimeout 設定超時
func WithTimeout(timeout time.Duration) HTTPClientOption {
	return func(c *HTTPClient) {
		c.client.Timeout = timeout
	}
}

// WithInsecureSkipVerify 跳過 TLS 驗證（僅限開發環境）
func WithInsecureSkipVerify() HTTPClientOption {
	return func(c *HTTPClient) {
		transport := c.client.Transport.(*http.Transport)
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
}

// NewHTTPClient 建立新的 HTTP 客戶端
func NewHTTPClient(opts ...HTTPClientOption) *HTTPClient {
	c := &HTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Do 執行 HTTP 請求
func (c *HTTPClient) Do(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	// 設定認證
	c.setAuth(req)

	return c.client.Do(req)
}

// Get 執行 GET 請求
func (c *HTTPClient) Get(ctx context.Context, path string) (*http.Response, error) {
	return c.Do(ctx, http.MethodGet, path, nil)
}

// Post 執行 POST 請求
func (c *HTTPClient) Post(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	return c.Do(ctx, http.MethodPost, path, body)
}

// setAuth 設定認證標頭
func (c *HTTPClient) setAuth(req *http.Request) {
	switch c.auth.Type {
	case "basic":
		req.SetBasicAuth(c.auth.Username, c.auth.Password)
	case "bearer":
		req.Header.Set("Authorization", "Bearer "+c.auth.Token)
	}
}

// Close 關閉客戶端
func (c *HTTPClient) Close() {
	c.client.CloseIdleConnections()
}
