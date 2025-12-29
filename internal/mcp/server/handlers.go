package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/vincent119/victorialogs-mcp/internal/util"
	"github.com/vincent119/victorialogs-mcp/internal/victorialogs"
)

// handleQuery handles vlogs-query request
func (s *MCPServer) handleQuery(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid request parameters"), nil
	}

	query, err := RequireString(args, "query")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	limit := GetInt(args, "limit", 1000)
	if limit > s.vlClient.GetMaxResults() {
		limit = s.vlClient.GetMaxResults()
	}

	// Parse time parameters
	var startTime, endTime *time.Time

	if start := GetString(args, "start", ""); start != "" {
		t, err := util.ParseTime(start)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid start time: %v", err)), nil
		}
		startTime = &t
	}

	if end := GetString(args, "end", ""); end != "" {
		t, err := util.ParseTime(end)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid end time: %v", err)), nil
		}
		endTime = &t
	}

	// Execute query
	result, err := s.vlClient.Query(ctx, victorialogs.QueryParams{
		Query: query,
		Start: startTime,
		End:   endTime,
		Limit: limit,
	})

	if err != nil {
		s.policyManager.RecordFailure()
		return mcp.NewToolResultError(fmt.Sprintf("query failed: %v", err)), nil
	}

	s.policyManager.RecordSuccess()

	// Format result
	output := formatQueryResult(result)
	return mcp.NewToolResultText(output), nil
}

// handleStats handles vlogs-stats request
func (s *MCPServer) handleStats(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid request parameters"), nil
	}

	start, err := RequireString(args, "start")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	startTime, err := util.ParseTime(start)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid start time: %v", err)), nil
	}

	query := GetString(args, "query", "")
	end := GetString(args, "end", "")

	var endTime *time.Time
	if end != "" {
		t, err := util.ParseTime(end)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid end time: %v", err)), nil
		}
		endTime = &t
	}

	result, err := s.vlClient.Stats(ctx, victorialogs.StatsParams{
		Query: query,
		Start: startTime,
		End:   endTime,
	})

	if err != nil {
		s.policyManager.RecordFailure()
		return mcp.NewToolResultError(fmt.Sprintf("stats query failed: %v", err)), nil
	}

	s.policyManager.RecordSuccess()

	// Format result
	output, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(output)), nil
}

// handleSchema handles vlogs-schema request
func (s *MCPServer) handleSchema(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid request parameters"), nil
	}

	schemaType, err := RequireString(args, "type")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	query := GetString(args, "query", "")
	field := GetString(args, "field", "")
	limit := GetInt(args, "limit", 100)

	result, err := s.vlClient.Schema(ctx, victorialogs.SchemaParams{
		Type:  schemaType,
		Query: query,
		Field: field,
		Limit: limit,
	})

	if err != nil {
		s.policyManager.RecordFailure()
		return mcp.NewToolResultError(fmt.Sprintf("schema query failed: %v", err)), nil
	}

	s.policyManager.RecordSuccess()

	// Format result
	output, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(output)), nil
}

// handleHealth handles vlogs-health request
func (s *MCPServer) handleHealth(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result, err := s.vlClient.Health(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("health check failed: %v", err)), nil
	}

	output, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(output)), nil
}

// formatQueryResult formats query result
func formatQueryResult(result *victorialogs.QueryResponse) string {
	var output string

	output += fmt.Sprintf("Found %d log entries", result.Total)
	if result.Truncated {
		output += " (results truncated)"
	}
	output += "\n\n"

	for i, entry := range result.Entries {
		output += fmt.Sprintf("--- [%d] %s ---\n", i+1, entry.Time.Format(time.RFC3339))
		if entry.Stream != "" {
			output += fmt.Sprintf("Stream: %s\n", entry.Stream)
		}
		output += fmt.Sprintf("Message: %s\n", entry.Message)

		if len(entry.Fields) > 0 {
			output += "Fields:\n"
			for k, v := range entry.Fields {
				output += fmt.Sprintf("  %s: %v\n", k, v)
			}
		}
		output += "\n"
	}

	return output
}
