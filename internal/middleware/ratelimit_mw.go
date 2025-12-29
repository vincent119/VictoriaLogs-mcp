// Package middleware 提供 MCP Tool 中介層
package middleware

import (
	"context"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/vincent119/victorialogs-mcp/internal/policy"
)

// RateLimitMiddleware Rate Limit 中介層
type RateLimitMiddleware struct {
	limiter *policy.RateLimiter
}

// NewRateLimitMiddleware 建立 Rate Limit 中介層
func NewRateLimitMiddleware(cfg policy.RateLimitConfig) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limiter: policy.NewRateLimiter(cfg),
	}
}

// Handler 回傳中介層處理函數
func (m *RateLimitMiddleware) Handler() ToolMiddleware {
	return func(next ToolHandler) ToolHandler {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// 使用請求的工具名稱作為 rate limit key
			key := request.Params.Name

			if err := m.limiter.Allow(key); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			return next(ctx, request)
		}
	}
}

// GetRemaining 取得剩餘請求次數
func (m *RateLimitMiddleware) GetRemaining(key string) int {
	return m.limiter.GetRemaining(key)
}

// StartCleanupRoutine 啟動清理協程
func (m *RateLimitMiddleware) StartCleanupRoutine(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				m.limiter.Cleanup()
			}
		}
	}()
}
