// Package server 提供 MCP Server 核心功能
package server

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/vincent119/victorialogs-mcp/internal/config"
	"github.com/vincent119/victorialogs-mcp/internal/middleware"
	"github.com/vincent119/victorialogs-mcp/internal/policy"
	"github.com/vincent119/victorialogs-mcp/internal/victorialogs"
	"github.com/vincent119/victorialogs-mcp/pkg/version"
	"github.com/vincent119/zlogger"
)

// MCPServer VictoriaLogs MCP Server
type MCPServer struct {
	server        *server.MCPServer
	vlClient      *victorialogs.Client
	policyManager *policy.Manager
	middlewares   []middleware.ToolMiddleware
	cfg           *config.Config
}

// New 建立新的 MCP Server
func New(cfg *config.Config, vlClient *victorialogs.Client, policyMgr *policy.Manager) *MCPServer {
	s := &MCPServer{
		cfg:           cfg,
		vlClient:      vlClient,
		policyManager: policyMgr,
		middlewares:   make([]middleware.ToolMiddleware, 0),
	}

	// 建立 MCP Server
	s.server = server.NewMCPServer(
		cfg.Server.Name,
		version.Short(),
		server.WithToolCapabilities(true),
		server.WithRecovery(),
	)

	// 建立 Middleware
	s.setupMiddlewares(cfg)

	// 註冊 Tools
	s.registerTools()

	return s
}

// setupMiddlewares 設定中介層
func (s *MCPServer) setupMiddlewares(cfg *config.Config) {
	// Rate Limit
	if cfg.Policy.RateLimit.Enabled {
		rateLimitMw := middleware.NewRateLimitMiddleware(policy.RateLimitConfig{
			Enabled:           cfg.Policy.RateLimit.Enabled,
			RequestsPerMinute: cfg.Policy.RateLimit.RequestsPerMinute,
		})
		s.middlewares = append(s.middlewares, rateLimitMw.Handler())
	}

	// Circuit Breaker
	if cfg.Policy.CircuitBreaker.Enabled {
		cbMw := middleware.NewCircuitBreakerMiddleware(policy.CircuitBreakerConfig{
			Enabled:        cfg.Policy.CircuitBreaker.Enabled,
			ErrorThreshold: cfg.Policy.CircuitBreaker.ErrorThreshold,
			Timeout:        cfg.Policy.CircuitBreaker.Timeout.String(),
		})
		s.middlewares = append(s.middlewares, cbMw.Handler())
	}

	// Audit
	auditMw := middleware.NewAuditMiddleware(middleware.AuditConfig{
		Enabled: true,
	})
	s.middlewares = append(s.middlewares, auditMw.Handler())

	// Redact（放在最後，處理輸出）
	// 使用預設的 Redact 規則
	redactMw := middleware.NewRedactMiddleware(policy.RedactConfig{
		Enabled: true,
	})
	s.middlewares = append(s.middlewares, redactMw.Handler())
}

// registerTools 註冊所有 Tools
func (s *MCPServer) registerTools() {
	// vlogs-query
	s.server.AddTool(
		mcp.NewTool("vlogs-query",
			mcp.WithDescription("Execute LogsQL query against VictoriaLogs. Returns matching log entries."),
			mcp.WithString("query",
				mcp.Required(),
				mcp.Description("LogsQL query string (e.g., 'error' or '_stream:{app=\"myapp\"}'"),
			),
			mcp.WithNumber("limit",
				mcp.Description("Maximum number of log entries to return (default: 1000, max: 5000)"),
			),
			mcp.WithString("start",
				mcp.Description("Start time - RFC3339 format or relative time like '5m', '1h', '24h'"),
			),
			mcp.WithString("end",
				mcp.Description("End time - RFC3339 format or relative time (default: now)"),
			),
		),
		s.wrapHandler(s.handleQuery),
	)

	// vlogs-stats
	s.server.AddTool(
		mcp.NewTool("vlogs-stats",
			mcp.WithDescription("Get log statistics (hit counts) over a time range."),
			mcp.WithString("query",
				mcp.Description("Optional LogsQL filter query"),
			),
			mcp.WithString("start",
				mcp.Required(),
				mcp.Description("Start time - RFC3339 format or relative time"),
			),
			mcp.WithString("end",
				mcp.Description("End time - RFC3339 format or relative time (default: now)"),
			),
		),
		s.wrapHandler(s.handleStats),
	)

	// vlogs-schema
	s.server.AddTool(
		mcp.NewTool("vlogs-schema",
			mcp.WithDescription("Explore available log streams, field names, or field values."),
			mcp.WithString("type",
				mcp.Required(),
				mcp.Description("Type of schema info to retrieve"),
				mcp.Enum("streams", "fields", "values"),
			),
			mcp.WithString("query",
				mcp.Description("Optional LogsQL filter query"),
			),
			mcp.WithString("field",
				mcp.Description("Field name for 'values' type query"),
			),
			mcp.WithNumber("limit",
				mcp.Description("Maximum number of results to return"),
			),
		),
		s.wrapHandler(s.handleSchema),
	)

	// vlogs-health
	s.server.AddTool(
		mcp.NewTool("vlogs-health",
			mcp.WithDescription("Check VictoriaLogs server health status."),
		),
		s.wrapHandler(s.handleHealth),
	)

	zlogger.Info("MCP Tools registered",
		zlogger.Int("count", 4),
		zlogger.String("tools", "vlogs-query, vlogs-stats, vlogs-schema, vlogs-health"),
	)
}

// wrapHandler 包裝 handler 並套用中介層
func (s *MCPServer) wrapHandler(handler middleware.ToolHandler) server.ToolHandlerFunc {
	// 串接所有中介層
	wrapped := middleware.Chain(s.middlewares...)(handler)

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return wrapped(ctx, request)
	}
}

// GetServer 取得底層 MCP Server
func (s *MCPServer) GetServer() *server.MCPServer {
	return s.server
}

// Close 關閉 Server
func (s *MCPServer) Close() error {
	if s.vlClient != nil {
		s.vlClient.Close()
	}
	if s.policyManager != nil {
		s.policyManager.Close()
	}
	return nil
}

// RequireString 從參數取得必填字串
func RequireString(args map[string]interface{}, key string) (string, error) {
	v, ok := args[key]
	if !ok {
		return "", fmt.Errorf("missing required parameter: %s", key)
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("parameter %s must be a string", key)
	}
	return s, nil
}

// GetString 從參數取得選填字串
func GetString(args map[string]interface{}, key, defaultValue string) string {
	v, ok := args[key]
	if !ok {
		return defaultValue
	}
	s, ok := v.(string)
	if !ok {
		return defaultValue
	}
	return s
}

// GetInt 從參數取得選填整數
func GetInt(args map[string]interface{}, key string, defaultValue int) int {
	v, ok := args[key]
	if !ok {
		return defaultValue
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	default:
		return defaultValue
	}
}
