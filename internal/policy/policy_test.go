package policy

import (
	"testing"
	"time"
)

func TestRateLimiter_Allow(t *testing.T) {
	cfg := RateLimitConfig{
		Enabled:           true,
		RequestsPerMinute: 5,
	}
	limiter := NewRateLimiter(cfg)

	// First 5 requests should be allowed
	for i := 0; i < 5; i++ {
		if err := limiter.Allow("test-key"); err != nil {
			t.Errorf("Request %d should be allowed, got error: %v", i+1, err)
		}
	}

	// 6th request should be denied
	if err := limiter.Allow("test-key"); err != ErrRateLimitExceeded {
		t.Errorf("Request 6 should be rate limited, got: %v", err)
	}

	// Different key should be allowed
	if err := limiter.Allow("other-key"); err != nil {
		t.Errorf("Different key should be allowed, got error: %v", err)
	}
}

func TestRateLimiter_Disabled(t *testing.T) {
	cfg := RateLimitConfig{
		Enabled:           false,
		RequestsPerMinute: 1,
	}
	limiter := NewRateLimiter(cfg)

	// All requests should be allowed when disabled
	for i := 0; i < 10; i++ {
		if err := limiter.Allow("test-key"); err != nil {
			t.Errorf("Request should be allowed when disabled, got error: %v", err)
		}
	}
}

func TestRateLimiter_GetRemaining(t *testing.T) {
	cfg := RateLimitConfig{
		Enabled:           true,
		RequestsPerMinute: 10,
	}
	limiter := NewRateLimiter(cfg)

	// Initially should have full quota
	if remaining := limiter.GetRemaining("test-key"); remaining != 10 {
		t.Errorf("Expected 10 remaining, got %d", remaining)
	}

	// After 3 requests, should have 7 remaining
	for i := 0; i < 3; i++ {
		limiter.Allow("test-key")
	}
	if remaining := limiter.GetRemaining("test-key"); remaining != 7 {
		t.Errorf("Expected 7 remaining, got %d", remaining)
	}
}

func TestRateLimiter_Reset(t *testing.T) {
	cfg := RateLimitConfig{
		Enabled:           true,
		RequestsPerMinute: 2,
	}
	limiter := NewRateLimiter(cfg)

	// Use up quota
	limiter.Allow("test-key")
	limiter.Allow("test-key")

	// Reset
	limiter.Reset("test-key")

	// Should be allowed again
	if err := limiter.Allow("test-key"); err != nil {
		t.Errorf("After reset should be allowed, got error: %v", err)
	}
}

func TestCircuitBreaker_Allow(t *testing.T) {
	cfg := CircuitBreakerConfig{
		Enabled:        true,
		ErrorThreshold: 3,
		Timeout:        "100ms",
	}
	cb := NewCircuitBreaker(cfg)

	// Initially should be closed and allow requests
	if err := cb.Allow(); err != nil {
		t.Errorf("Initially should allow, got error: %v", err)
	}

	// Record failures to open circuit
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}

	// Should be open now
	if err := cb.Allow(); err != ErrCircuitOpen {
		t.Errorf("After failures should be open, got: %v", err)
	}

	// Wait for timeout
	time.Sleep(150 * time.Millisecond)

	// Should allow (half-open)
	if err := cb.Allow(); err != nil {
		t.Errorf("After timeout should allow (half-open), got error: %v", err)
	}

	// Record success to close
	cb.RecordSuccess()

	if state := cb.GetStateString(); state != "closed" {
		t.Errorf("After success should be closed, got: %s", state)
	}
}

func TestCircuitBreaker_Disabled(t *testing.T) {
	cfg := CircuitBreakerConfig{
		Enabled:        false,
		ErrorThreshold: 1,
		Timeout:        "1s",
	}
	cb := NewCircuitBreaker(cfg)

	// Should always allow when disabled
	for i := 0; i < 10; i++ {
		cb.RecordFailure()
		if err := cb.Allow(); err != nil {
			t.Errorf("Should always allow when disabled, got error: %v", err)
		}
	}
}

func TestAllowlist_Check(t *testing.T) {
	cfg := AllowlistConfig{
		Enabled: true,
		Streams: []string{"app/*", "kubernetes/**"},
		Deny:    []string{"secret/*"},
	}
	allowlist := NewAllowlist(cfg)

	tests := []struct {
		stream  string
		allowed bool
	}{
		{"app/logs", true},
		{"app/metrics", true},
		{"kubernetes/pod", true},
		{"kubernetes/node/events", true},
		{"secret/keys", false}, // explicitly denied
		{"unknown/stream", false},
	}

	for _, tt := range tests {
		err := allowlist.Check(tt.stream)
		if tt.allowed && err != nil {
			t.Errorf("Stream %s should be allowed, got error: %v", tt.stream, err)
		}
		if !tt.allowed && err == nil {
			t.Errorf("Stream %s should not be allowed", tt.stream)
		}
	}
}

func TestAllowlist_Disabled(t *testing.T) {
	cfg := AllowlistConfig{
		Enabled: false,
		Streams: []string{"allowed/*"},
	}
	allowlist := NewAllowlist(cfg)

	// Should allow any stream when disabled
	if err := allowlist.Check("any/stream"); err != nil {
		t.Errorf("Should allow any stream when disabled, got error: %v", err)
	}
}
