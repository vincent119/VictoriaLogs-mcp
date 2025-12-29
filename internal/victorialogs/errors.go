// Package victorialogs 提供 VictoriaLogs API 客戶端
package victorialogs

import (
	"fmt"
)

// 錯誤類型定義
var (
	// ErrConnection 連線錯誤
	ErrConnection = fmt.Errorf("VictoriaLogs connection error")

	// ErrQueryFailed 查詢失敗
	ErrQueryFailed = fmt.Errorf("query execution failed")

	// ErrInvalidQuery 無效的查詢
	ErrInvalidQuery = fmt.Errorf("invalid LogsQL query")

	// ErrTimeout 查詢超時
	ErrTimeout = fmt.Errorf("query timeout")

	// ErrUnauthorized 認證失敗
	ErrUnauthorized = fmt.Errorf("authentication failed")

	// ErrTooManyRequests 請求過於頻繁
	ErrTooManyRequests = fmt.Errorf("too many requests")
)

// APIError VictoriaLogs API 錯誤
type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Query      string `json:"query,omitempty"`
}

// Error 實作 error 介面
func (e *APIError) Error() string {
	if e.Query != "" {
		return fmt.Sprintf("VictoriaLogs API error (HTTP %d): %s, Query: %s", e.StatusCode, e.Message, e.Query)
	}
	return fmt.Sprintf("VictoriaLogs API error (HTTP %d): %s", e.StatusCode, e.Message)
}

// IsConnectionError 檢查是否為連線錯誤
func IsConnectionError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == 0 || apiErr.StatusCode >= 500
	}
	return false
}

// IsAuthError 檢查是否為認證錯誤
func IsAuthError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == 401 || apiErr.StatusCode == 403
	}
	return false
}

// IsRateLimitError 檢查是否為速率限制錯誤
func IsRateLimitError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == 429
	}
	return false
}

// NewAPIError 建立 API 錯誤
func NewAPIError(statusCode int, message, query string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		Query:      query,
	}
}
