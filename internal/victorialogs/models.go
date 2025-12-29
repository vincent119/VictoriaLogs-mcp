package victorialogs

import (
	"time"
)

// LogEntry 日誌條目
type LogEntry struct {
	Time    time.Time              `json:"_time"`
	Message string                 `json:"_msg"`
	Stream  string                 `json:"_stream,omitempty"`
	Fields  map[string]interface{} `json:"-"`
}

// QueryResponse 查詢回應
type QueryResponse struct {
	Entries   []LogEntry `json:"entries"`
	Total     int        `json:"total"`
	Truncated bool       `json:"truncated"`
}

// StatsResponse 統計回應
type StatsResponse struct {
	Hits []HitEntry `json:"hits"`
}

// HitEntry 統計條目
type HitEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Count     int64     `json:"count"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// StreamInfo Stream 資訊
type StreamInfo struct {
	Stream string            `json:"_stream"`
	Labels map[string]string `json:"labels,omitempty"`
}

// StreamsResponse Streams 查詢回應
type StreamsResponse struct {
	Streams []StreamInfo `json:"streams"`
}

// FieldInfo 欄位資訊
type FieldInfo struct {
	Name  string `json:"name"`
	Hits  int64  `json:"hits"`
}

// FieldsResponse 欄位查詢回應
type FieldsResponse struct {
	Fields []FieldInfo `json:"fields"`
}

// FieldValuesResponse 欄位值查詢回應
type FieldValuesResponse struct {
	Values []string `json:"values"`
}

// HealthResponse 健康檢查回應
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version,omitempty"`
}

// QueryParams 查詢參數
type QueryParams struct {
	Query string     `json:"query"`
	Start *time.Time `json:"start,omitempty"`
	End   *time.Time `json:"end,omitempty"`
	Limit int        `json:"limit,omitempty"`
}

// StatsParams 統計參數
type StatsParams struct {
	Query string     `json:"query,omitempty"`
	Start time.Time  `json:"start"`
	End   *time.Time `json:"end,omitempty"`
	Step  string     `json:"step,omitempty"`
}

// SchemaParams Schema 查詢參數
type SchemaParams struct {
	Type  string `json:"type"` // streams | fields | values
	Query string `json:"query,omitempty"`
	Field string `json:"field,omitempty"` // 用於 values 查詢
	Limit int    `json:"limit,omitempty"`
}
