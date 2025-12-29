// Package schema 提供 MCP Tool Schema 定義
package schema

import "github.com/mark3labs/mcp-go/mcp"

// VLogsQuery vlogs-query Tool 定義
var VLogsQuery = mcp.NewTool("vlogs-query",
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
)

// VLogsStats vlogs-stats Tool 定義
var VLogsStats = mcp.NewTool("vlogs-stats",
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
)

// VLogsSchema vlogs-schema Tool 定義
var VLogsSchema = mcp.NewTool("vlogs-schema",
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
)

// VLogsTail vlogs-tail Tool 定義
var VLogsTail = mcp.NewTool("vlogs-tail",
	mcp.WithDescription("Stream live log entries matching the query. Returns a limited number of recent entries."),
	mcp.WithString("query",
		mcp.Required(),
		mcp.Description("LogsQL query string to filter logs"),
	),
	mcp.WithNumber("limit",
		mcp.Description("Maximum number of log entries to return (default: 100)"),
	),
	mcp.WithNumber("timeout",
		mcp.Description("Maximum time in seconds to wait for logs (default: 5)"),
	),
)

// VLogsExplain vlogs-explain Tool 定義
var VLogsExplain = mcp.NewTool("vlogs-explain",
	mcp.WithDescription("Explain the execution plan of a LogsQL query."),
	mcp.WithString("query",
		mcp.Required(),
		mcp.Description("LogsQL query string to explain"),
	),
)

// VLogsHealth vlogs-health Tool 定義
var VLogsHealth = mcp.NewTool("vlogs-health",
	mcp.WithDescription("Check VictoriaLogs server health status."),
)

// AllTools 所有 Tool 定義
var AllTools = []mcp.Tool{
	VLogsQuery,
	VLogsStats,
	VLogsSchema,
	VLogsTail,
	VLogsExplain,
	VLogsHealth,
}
