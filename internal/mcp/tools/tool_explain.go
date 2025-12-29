package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/vincent119/victorialogs-mcp/internal/victorialogs"
)

// ExplainHandler vlogs-explain tool handler
type ExplainHandler struct {
	client *victorialogs.Client
}

// NewExplainHandler creates explain handler
func NewExplainHandler(client *victorialogs.Client) *ExplainHandler {
	return &ExplainHandler{client: client}
}

// Handle handles vlogs-explain request
func (h *ExplainHandler) Handle(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid request parameters"), nil
	}

	query, ok := args["query"].(string)
	if !ok || query == "" {
		return mcp.NewToolResultError("missing required parameter: query"), nil
	}

	// Execute query analysis
	// Note: VictoriaLogs may not have a dedicated explain endpoint
	// We provide syntax analysis and recommendations

	// Build explanation result
	result := fmt.Sprintf(`LogsQL Query Analysis
=====================

Original Query:
%s

Description:
- This query searches for matching logs in VictoriaLogs
- Uses LogsQL syntax for filtering and aggregation
- Recommend adding time range for better query performance

LogsQL Syntax Tips:
- Use _stream:{key="value"} to filter specific streams
- Use AND, OR, NOT for logical combinations
- Use | stats for statistical aggregation
- Use | fields to select specific fields

Documentation: https://docs.victoriametrics.com/victorialogs/logsql/
`, query)

	return mcp.NewToolResultText(result), nil
}

// Explain executes query explanation
func (h *ExplainHandler) Explain(_ context.Context, query string) (string, error) {
	if query == "" {
		return "", victorialogs.ErrInvalidQuery
	}

	// Basic syntax analysis
	analysis := analyzeQuery(query)
	return analysis, nil
}

// analyzeQuery analyzes LogsQL query
func analyzeQuery(query string) string {
	var analysis string

	analysis += fmt.Sprintf("Query: %s\n\n", query)

	// Check common patterns
	if len(query) < 3 {
		analysis += "⚠️ Warning: Query is too short, may match many results\n"
	}

	// Check if stream filter is used
	if contains(query, "_stream:") {
		analysis += "✅ Using _stream filter helps narrow search scope\n"
	} else {
		analysis += "💡 Tip: Adding _stream filter can improve query efficiency\n"
	}

	// Check if time filter is used
	if contains(query, "_time:") {
		analysis += "✅ Using time filter\n"
	}

	// Check if stats operation is used
	if contains(query, "| stats") || contains(query, "|stats") {
		analysis += "📊 Contains statistical aggregation\n"
	}

	return analysis
}

// contains checks if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
