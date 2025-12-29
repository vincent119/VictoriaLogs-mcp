package victorialogs

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Streams 查詢日誌 Streams
func (c *Client) Streams(ctx context.Context, query string, limit int) (*StreamsResponse, error) {
	params := url.Values{}

	if query != "" {
		params.Set("query", query)
	}

	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	body, err := c.doRequest(ctx, "GET", "/select/logsql/streams", params)
	if err != nil {
		return nil, err
	}

	streams, err := parseStreamsNDJSON(body)
	if err != nil {
		return nil, fmt.Errorf("parse failed: %w", err)
	}

	return &StreamsResponse{Streams: streams}, nil
}

// FieldNames 查詢欄位名稱
func (c *Client) FieldNames(ctx context.Context, query string, limit int) (*FieldsResponse, error) {
	params := url.Values{}

	if query != "" {
		params.Set("query", query)
	}

	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	body, err := c.doRequest(ctx, "GET", "/select/logsql/field_names", params)
	if err != nil {
		return nil, err
	}

	fields, err := parseFieldsNDJSON(body)
	if err != nil {
		return nil, fmt.Errorf("parse failed: %w", err)
	}

	return &FieldsResponse{Fields: fields}, nil
}

// FieldValues 查詢欄位值
func (c *Client) FieldValues(ctx context.Context, field, query string, limit int) (*FieldValuesResponse, error) {
	if field == "" {
		return nil, fmt.Errorf("field is required")
	}

	params := url.Values{}
	params.Set("field", field)

	if query != "" {
		params.Set("query", query)
	}

	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	body, err := c.doRequest(ctx, "GET", "/select/logsql/field_values", params)
	if err != nil {
		return nil, err
	}

	values, err := parseFieldValuesNDJSON(body)
	if err != nil {
		return nil, fmt.Errorf("parse failed: %w", err)
	}

	return &FieldValuesResponse{Values: values}, nil
}

// Schema 統一的 Schema 查詢介面
func (c *Client) Schema(ctx context.Context, params SchemaParams) (interface{}, error) {
	switch params.Type {
	case "streams":
		return c.Streams(ctx, params.Query, params.Limit)
	case "fields":
		return c.FieldNames(ctx, params.Query, params.Limit)
	case "values":
		return c.FieldValues(ctx, params.Field, params.Query, params.Limit)
	default:
		return nil, fmt.Errorf("unsupported schema type: %s", params.Type)
	}
}

// parseStreamsNDJSON 解析 streams NDJSON
func parseStreamsNDJSON(data []byte) ([]StreamInfo, error) {
	var streams []StreamInfo
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

		info := StreamInfo{
			Labels: make(map[string]string),
		}

		if stream, ok := raw["_stream"]; ok {
			if s, ok := stream.(string); ok {
				info.Stream = s
			}
		}

		for k, v := range raw {
			if k != "_stream" {
				if s, ok := v.(string); ok {
					info.Labels[k] = s
				}
			}
		}

		streams = append(streams, info)
	}

	return streams, scanner.Err()
}

// parseFieldsNDJSON 解析 fields NDJSON
func parseFieldsNDJSON(data []byte) ([]FieldInfo, error) {
	var fields []FieldInfo
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

		for name, hits := range raw {
			info := FieldInfo{Name: name}
			if h, ok := hits.(float64); ok {
				info.Hits = int64(h)
			}
			fields = append(fields, info)
		}
	}

	return fields, scanner.Err()
}

// parseFieldValuesNDJSON 解析 field values NDJSON
func parseFieldValuesNDJSON(data []byte) ([]string, error) {
	var values []string
	scanner := bufio.NewScanner(string2reader(data))

	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			// 移除引號
			line = strings.Trim(line, "\"")
			values = append(values, line)
		}
	}

	return values, scanner.Err()
}
