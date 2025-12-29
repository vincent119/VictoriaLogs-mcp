package victorialogs

import (
	"bufio"
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/vincent119/victorialogs-mcp/internal/util"
)

// Stats 查詢日誌統計
func (c *Client) Stats(ctx context.Context, params StatsParams) (*StatsResponse, error) {
	query := url.Values{}

	if params.Query != "" {
		query.Set("query", params.Query)
	}

	query.Set("start", util.FormatTime(params.Start))

	if params.End != nil {
		query.Set("end", util.FormatTime(*params.End))
	}

	if params.Step != "" {
		query.Set("step", params.Step)
	}

	body, err := c.doRequest(ctx, "GET", "/select/logsql/hits", query)
	if err != nil {
		return nil, err
	}

	var response StatsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		// 嘗試解析為 NDJSON 格式
		hits, parseErr := parseHitsNDJSON(body)
		if parseErr != nil {
			return nil, parseErr
		}
		response.Hits = hits
	}

	return &response, nil
}

// StatsQuery 執行統計查詢
func (c *Client) StatsQuery(ctx context.Context, queryStr string, start string, end string) (*StatsResponse, error) {
	startTime, err := util.ParseTime(start)
	if err != nil {
		return nil, err
	}

	var endTime *time.Time
	if end != "" {
		t, err := util.ParseTime(end)
		if err != nil {
			return nil, err
		}
		endTime = &t
	}

	return c.Stats(ctx, StatsParams{
		Query: queryStr,
		Start: startTime,
		End:   endTime,
	})
}

// parseHitsNDJSON 解析 hits NDJSON 格式
func parseHitsNDJSON(data []byte) ([]HitEntry, error) {
	var hits []HitEntry
	scanner := bufio.NewScanner(string2reader(data))

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var raw map[string]interface{}
		if err := json.Unmarshal(line, &raw); err != nil {
			continue
		}

		entry := HitEntry{
			Fields: make(map[string]interface{}),
		}

		if h, ok := raw["hits"]; ok {
			if count, ok := h.(float64); ok {
				entry.Count = int64(count)
			}
		}

		for k, v := range raw {
			if k != "hits" {
				entry.Fields[k] = v
			}
		}

		hits = append(hits, entry)
	}

	return hits, scanner.Err()
}
