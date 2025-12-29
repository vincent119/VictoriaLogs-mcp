package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/vincent119/zlogger"
)

// AuditMiddleware audit logging middleware
type AuditMiddleware struct {
	enabled bool
}

// AuditConfig audit configuration
type AuditConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// AuditEntry audit log entry
type AuditEntry struct {
	Timestamp     time.Time         `json:"timestamp"`
	ToolName      string            `json:"tool_name"`
	Duration      time.Duration     `json:"duration_ms"`
	Success       bool              `json:"success"`
	Error         string            `json:"error,omitempty"`
	RequestID     string            `json:"request_id,omitempty"`
	ParamsSummary map[string]string `json:"params_summary,omitempty"`
}

// NewAuditMiddleware creates audit middleware
func NewAuditMiddleware(cfg AuditConfig) *AuditMiddleware {
	return &AuditMiddleware{
		enabled: cfg.Enabled,
	}
}

// Handler returns middleware handler function
func (m *AuditMiddleware) Handler() ToolMiddleware {
	return func(next ToolHandler) ToolHandler {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			if !m.enabled {
				return next(ctx, request)
			}

			start := time.Now()

			// Execute actual handler
			result, err := next(ctx, request)

			// Log audit entry
			entry := AuditEntry{
				Timestamp:     start,
				ToolName:      request.Params.Name,
				Duration:      time.Since(start),
				Success:       err == nil && (result == nil || !result.IsError),
				ParamsSummary: m.extractParamsSummary(request),
			}

			if err != nil {
				entry.Error = err.Error()
			} else if result != nil && result.IsError {
				entry.Error = "tool returned error"
			}

			m.logEntry(entry)

			return result, err
		}
	}
}

// extractParamsSummary extracts parameter summary (avoid logging sensitive info)
func (m *AuditMiddleware) extractParamsSummary(request mcp.CallToolRequest) map[string]string {
	summary := make(map[string]string)

	args := request.Params.Arguments
	if args == nil {
		return summary
	}

	// Convert args to map[string]interface{}
	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return summary
	}

	// Only log non-sensitive parameter summaries
	if query, ok := argsMap["query"]; ok {
		if q, ok := query.(string); ok {
			if len(q) > 50 {
				summary["query_preview"] = q[:50] + "..."
			} else {
				summary["query_preview"] = q
			}
		}
	}

	if start, ok := argsMap["start"]; ok {
		if s, ok := start.(string); ok {
			summary["start"] = s
		}
	}

	if end, ok := argsMap["end"]; ok {
		if e, ok := end.(string); ok {
			summary["end"] = e
		}
	}

	if limit, ok := argsMap["limit"]; ok {
		switch l := limit.(type) {
		case float64:
			summary["limit"] = fmt.Sprintf("%d", int(l))
		case int:
			summary["limit"] = fmt.Sprintf("%d", l)
		}
	}

	return summary
}

// logEntry logs audit entry
func (m *AuditMiddleware) logEntry(entry AuditEntry) {
	if entry.Success {
		zlogger.Info("MCP Tool call",
			zlogger.String("tool", entry.ToolName),
			zlogger.Int64("duration_ms", int64(entry.Duration.Milliseconds())),
			zlogger.Bool("success", entry.Success),
		)
	} else {
		zlogger.Warn("MCP Tool call failed",
			zlogger.String("tool", entry.ToolName),
			zlogger.Int64("duration_ms", int64(entry.Duration.Milliseconds())),
			zlogger.Bool("success", entry.Success),
			zlogger.String("error", entry.Error),
		)
	}
}
