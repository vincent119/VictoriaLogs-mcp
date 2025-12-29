package victorialogs

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/vincent119/victorialogs-mcp/internal/util"
)

// Query 執行 LogsQL 查詢
func (c *Client) Query(ctx context.Context, params QueryParams) (*QueryResponse, error) {
	if params.Query == "" {
		return nil, ErrInvalidQuery
	}

	query := url.Values{}
	query.Set("query", params.Query)

	if params.Start != nil {
		query.Set("start", util.FormatTime(*params.Start))
	}
	if params.End != nil {
		query.Set("end", util.FormatTime(*params.End))
	}

	limit := params.Limit
	if limit <= 0 || limit > c.maxResults {
		limit = c.maxResults
	}
	query.Set("limit", strconv.Itoa(limit))

	body, err := c.doRequest(ctx, "GET", "/select/logsql/query", query)
	if err != nil {
		if apiErr, ok := err.(*APIError); ok {
			apiErr.Query = params.Query
		}
		return nil, err
	}

	// VictoriaLogs 回傳 NDJSON 格式（每行一個 JSON 物件）
	entries, err := parseNDJSON(body)
	if err != nil {
		return nil, fmt.Errorf("解析查詢結果失敗: %w", err)
	}

	truncated := len(entries) >= limit

	return &QueryResponse{
		Entries:   entries,
		Total:     len(entries),
		Truncated: truncated,
	}, nil
}

// QueryWithTimeRange 使用時間範圍執行查詢
func (c *Client) QueryWithTimeRange(ctx context.Context, queryStr string, start, end time.Time, limit int) (*QueryResponse, error) {
	return c.Query(ctx, QueryParams{
		Query: queryStr,
		Start: &start,
		End:   &end,
		Limit: limit,
	})
}

// QueryRelative 使用相對時間執行查詢
func (c *Client) QueryRelative(ctx context.Context, queryStr, relativeTime string, limit int) (*QueryResponse, error) {
	start, err := util.ParseTime(relativeTime)
	if err != nil {
		return nil, fmt.Errorf("解析相對時間失敗: %w", err)
	}

	now := time.Now()
	return c.Query(ctx, QueryParams{
		Query: queryStr,
		Start: &start,
		End:   &now,
		Limit: limit,
	})
}

// parseNDJSON 解析 NDJSON 格式
func parseNDJSON(data []byte) ([]LogEntry, error) {
	var entries []LogEntry
	scanner := bufio.NewScanner(string2reader(data))

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var raw map[string]interface{}
		if err := json.Unmarshal(line, &raw); err != nil {
			continue // 跳過無效的行
		}

		entry := LogEntry{
			Fields: make(map[string]interface{}),
		}

		// 解析標準欄位
		if t, ok := raw["_time"]; ok {
			if ts, ok := t.(string); ok {
				if parsed, err := time.Parse(time.RFC3339Nano, ts); err == nil {
					entry.Time = parsed
				}
			}
		}

		if msg, ok := raw["_msg"]; ok {
			if s, ok := msg.(string); ok {
				entry.Message = s
			}
		}

		if stream, ok := raw["_stream"]; ok {
			if s, ok := stream.(string); ok {
				entry.Stream = s
			}
		}

		// 其他欄位放入 Fields
		for k, v := range raw {
			if k != "_time" && k != "_msg" && k != "_stream" {
				entry.Fields[k] = v
			}
		}

		entries = append(entries, entry)
	}

	return entries, scanner.Err()
}

// string2reader 將 byte slice 轉為 io.Reader
func string2reader(data []byte) *bufio.Reader {
	return bufio.NewReader(&byteReader{data: data})
}

type byteReader struct {
	data []byte
	pos  int
}

func (r *byteReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, fmt.Errorf("EOF")
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}
