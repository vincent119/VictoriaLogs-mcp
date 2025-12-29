package middleware

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/vincent119/victorialogs-mcp/internal/policy"
)

// CircuitBreakerMiddleware Circuit Breaker 中介層
type CircuitBreakerMiddleware struct {
	cb *policy.CircuitBreaker
}

// NewCircuitBreakerMiddleware 建立 Circuit Breaker 中介層
func NewCircuitBreakerMiddleware(cfg policy.CircuitBreakerConfig) *CircuitBreakerMiddleware {
	return &CircuitBreakerMiddleware{
		cb: policy.NewCircuitBreaker(cfg),
	}
}

// Handler 回傳中介層處理函數
func (m *CircuitBreakerMiddleware) Handler() ToolMiddleware {
	return func(next ToolHandler) ToolHandler {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// 檢查 Circuit Breaker 狀態
			if err := m.cb.Allow(); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			// 執行實際的處理
			result, err := next(ctx, request)

			// 根據結果更新 Circuit Breaker
			if err != nil || (result != nil && result.IsError) {
				m.cb.RecordFailure()
			} else {
				m.cb.RecordSuccess()
			}

			return result, err
		}
	}
}

// GetState 取得當前狀態
func (m *CircuitBreakerMiddleware) GetState() string {
	return m.cb.GetStateString()
}

// Reset 重置 Circuit Breaker
func (m *CircuitBreakerMiddleware) Reset() {
	m.cb.Reset()
}

// ToolHandler Tool 處理函數型別
type ToolHandler func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)

// ToolMiddleware Tool 中介層型別
type ToolMiddleware func(next ToolHandler) ToolHandler

// Chain 串接多個中介層
func Chain(middlewares ...ToolMiddleware) ToolMiddleware {
	return func(final ToolHandler) ToolHandler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}

// NoopMiddleware 空操作中介層（用於測試）
func NoopMiddleware() ToolMiddleware {
	return func(next ToolHandler) ToolHandler {
		return next
	}
}
