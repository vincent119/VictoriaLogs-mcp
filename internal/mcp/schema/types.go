// Package schema 提供 MCP 共用型別定義
package schema

import "time"

// QueryResult 通用查詢結果
type QueryResult struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Truncated bool        `json:"truncated,omitempty"`
	Total     int         `json:"total,omitempty"`
}

// LogEntry 日誌條目
type LogEntry struct {
	Time    time.Time              `json:"_time"`
	Message string                 `json:"_msg"`
	Stream  string                 `json:"_stream,omitempty"`
	Fields  map[string]interface{} `json:"fields,omitempty"`
}

// StatsEntry 統計條目
type StatsEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Count     int64                  `json:"count"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// StreamInfo Stream 資訊
type StreamInfo struct {
	Stream string            `json:"_stream"`
	Labels map[string]string `json:"labels,omitempty"`
}

// FieldInfo 欄位資訊
type FieldInfo struct {
	Name string `json:"name"`
	Hits int64  `json:"hits"`
}

// HealthStatus 健康狀態
type HealthStatus struct {
	Status  string `json:"status"`
	Version string `json:"version,omitempty"`
	Message string `json:"message,omitempty"`
}

// ExplainResult 執行計畫結果
type ExplainResult struct {
	Query       string `json:"query"`
	ParsedQuery string `json:"parsed_query,omitempty"`
	Plan        string `json:"plan,omitempty"`
}

// ToolNames Tool 名稱常量
const (
	ToolQuery   = "vlogs-query"
	ToolStats   = "vlogs-stats"
	ToolSchema  = "vlogs-schema"
	ToolTail    = "vlogs-tail"
	ToolExplain = "vlogs-explain"
	ToolHealth  = "vlogs-health"
)

// SchemaTypes Schema 查詢類型
const (
	SchemaTypeStreams = "streams"
	SchemaTypeFields  = "fields"
	SchemaTypeValues  = "values"
)
