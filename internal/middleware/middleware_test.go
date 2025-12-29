package middleware

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/vincent119/victorialogs-mcp/internal/policy"
)

// Helper to create a simple request
func newTestRequest(_ string) mcp.CallToolRequest {
	return mcp.CallToolRequest{}
}

func TestRateLimitMiddleware(t *testing.T) {
	cfg := policy.RateLimitConfig{
		Enabled:           true,
		RequestsPerMinute: 3,
	}
	mw := NewRateLimitMiddleware(cfg)

	handler := func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("success"), nil
	}

	wrapped := mw.Handler()(handler)
	ctx := context.Background()
	req := newTestRequest("test-tool")

	// First 3 requests should succeed
	for i := 0; i < 3; i++ {
		result, err := wrapped(ctx, req)
		if err != nil {
			t.Errorf("Request %d should succeed: %v", i+1, err)
		}
		if result.IsError {
			t.Errorf("Request %d should not be error", i+1)
		}
	}

	// 4th request should be rate limited
	result, err := wrapped(ctx, req)
	if err != nil {
		t.Errorf("Should return result, not error: %v", err)
	}
	if !result.IsError {
		t.Error("4th request should be rate limited")
	}
}

func TestRateLimitMiddleware_Disabled(t *testing.T) {
	cfg := policy.RateLimitConfig{
		Enabled:           false,
		RequestsPerMinute: 1,
	}
	mw := NewRateLimitMiddleware(cfg)

	handler := func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("success"), nil
	}

	wrapped := mw.Handler()(handler)
	ctx := context.Background()
	req := newTestRequest("test-tool")

	// All requests should succeed when disabled
	for i := 0; i < 10; i++ {
		result, err := wrapped(ctx, req)
		if err != nil {
			t.Errorf("Request should succeed when disabled: %v", err)
		}
		if result.IsError {
			t.Error("Request should not be error when disabled")
		}
	}
}

func TestCircuitBreakerMiddleware(t *testing.T) {
	cfg := policy.CircuitBreakerConfig{
		Enabled:        true,
		ErrorThreshold: 2,
		Timeout:        "100ms",
	}
	mw := NewCircuitBreakerMiddleware(cfg)

	handler := func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("success"), nil
	}

	wrapped := mw.Handler()(handler)
	ctx := context.Background()
	req := newTestRequest("test-tool")

	// First request should succeed
	result, err := wrapped(ctx, req)
	if err != nil {
		t.Errorf("First request should succeed: %v", err)
	}
	if result.IsError {
		t.Error("First request should not be error")
	}
}

func TestAuditMiddleware(t *testing.T) {
	cfg := AuditConfig{
		Enabled: true,
	}
	mw := NewAuditMiddleware(cfg)

	handler := func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("success"), nil
	}

	wrapped := mw.Handler()(handler)
	ctx := context.Background()
	req := newTestRequest("test-tool")

	// Should log audit entry (no error)
	result, err := wrapped(ctx, req)
	if err != nil {
		t.Errorf("Audit middleware should not error: %v", err)
	}
	if result.IsError {
		t.Error("Result should not be error")
	}
}

func TestRedactMiddleware(t *testing.T) {
	cfg := policy.RedactConfig{
		Enabled: true,
	}
	mw := NewRedactMiddleware(cfg)

	// Handler returns text with sensitive data
	handler := func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("password=secret123"), nil
	}

	wrapped := mw.Handler()(handler)
	ctx := context.Background()
	req := newTestRequest("test-tool")

	result, err := wrapped(ctx, req)
	if err != nil {
		t.Errorf("Redact middleware should not error: %v", err)
	}
	if result == nil {
		t.Fatal("Result should not be nil")
	}
}

func TestChain(t *testing.T) {
	calls := []string{}

	mw1 := func(next ToolHandler) ToolHandler {
		return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			calls = append(calls, "mw1-before")
			result, err := next(ctx, req)
			calls = append(calls, "mw1-after")
			return result, err
		}
	}

	mw2 := func(next ToolHandler) ToolHandler {
		return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			calls = append(calls, "mw2-before")
			result, err := next(ctx, req)
			calls = append(calls, "mw2-after")
			return result, err
		}
	}

	handler := func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		calls = append(calls, "handler")
		return mcp.NewToolResultText("success"), nil
	}

	chained := Chain(mw1, mw2)(handler)
	_, _ = chained(context.Background(), newTestRequest("test"))

	expected := []string{"mw1-before", "mw2-before", "handler", "mw2-after", "mw1-after"}
	if len(calls) != len(expected) {
		t.Errorf("Expected %d calls, got %d", len(expected), len(calls))
	}

	for i, call := range calls {
		if call != expected[i] {
			t.Errorf("Call %d: expected %s, got %s", i, expected[i], call)
		}
	}
}

func TestNoopMiddleware(t *testing.T) {
	handler := func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("original"), nil
	}

	wrapped := NoopMiddleware()(handler)
	result, err := wrapped(context.Background(), newTestRequest("test"))

	if err != nil {
		t.Error("Noop should not error")
	}
	if result == nil {
		t.Error("Result should not be nil")
	}
}
