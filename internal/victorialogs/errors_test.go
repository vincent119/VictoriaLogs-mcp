package victorialogs

import "testing"

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		err      *APIError
		contains string
	}{
		{
			err:      &APIError{StatusCode: 404, Message: "not found"},
			contains: "404",
		},
		{
			err:      &APIError{StatusCode: 500, Message: "server error", Query: "test query"},
			contains: "test query",
		},
	}

	for _, tt := range tests {
		result := tt.err.Error()
		if !containsStr(result, tt.contains) {
			t.Errorf("Error message %q should contain %q", result, tt.contains)
		}
	}
}

func TestIsConnectionError(t *testing.T) {
	tests := []struct {
		err      error
		expected bool
	}{
		{&APIError{StatusCode: 0, Message: "connection refused"}, true},
		{&APIError{StatusCode: 500, Message: "internal error"}, true},
		{&APIError{StatusCode: 404, Message: "not found"}, false},
		{&APIError{StatusCode: 200, Message: "ok"}, false},
	}

	for _, tt := range tests {
		result := IsConnectionError(tt.err)
		if result != tt.expected {
			t.Errorf("IsConnectionError(%v) = %v, want %v", tt.err, result, tt.expected)
		}
	}
}

func TestIsAuthError(t *testing.T) {
	tests := []struct {
		err      error
		expected bool
	}{
		{&APIError{StatusCode: 401, Message: "unauthorized"}, true},
		{&APIError{StatusCode: 403, Message: "forbidden"}, true},
		{&APIError{StatusCode: 404, Message: "not found"}, false},
	}

	for _, tt := range tests {
		result := IsAuthError(tt.err)
		if result != tt.expected {
			t.Errorf("IsAuthError(%v) = %v, want %v", tt.err, result, tt.expected)
		}
	}
}

func TestIsRateLimitError(t *testing.T) {
	tests := []struct {
		err      error
		expected bool
	}{
		{&APIError{StatusCode: 429, Message: "too many requests"}, true},
		{&APIError{StatusCode: 500, Message: "server error"}, false},
	}

	for _, tt := range tests {
		result := IsRateLimitError(tt.err)
		if result != tt.expected {
			t.Errorf("IsRateLimitError(%v) = %v, want %v", tt.err, result, tt.expected)
		}
	}
}

func TestNewAPIError(t *testing.T) {
	err := NewAPIError(404, "not found", "select * from logs")

	if err.StatusCode != 404 {
		t.Errorf("Expected status 404, got %d", err.StatusCode)
	}
	if err.Message != "not found" {
		t.Errorf("Expected message 'not found', got %s", err.Message)
	}
	if err.Query != "select * from logs" {
		t.Errorf("Expected query, got %s", err.Query)
	}
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStrHelper(s, substr))
}

func containsStrHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
