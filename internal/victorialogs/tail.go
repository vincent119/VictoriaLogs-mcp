package victorialogs

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// TailCallback tail callback function
type TailCallback func(entry LogEntry) error

// Tail streams live logs (note: this is a blocking operation)
// This method will continue reading until context is cancelled or error occurs
func (c *Client) Tail(ctx context.Context, query string, callback TailCallback) error {
	if query == "" {
		return ErrInvalidQuery
	}

	params := url.Values{}
	params.Set("query", query)

	fullPath := "/select/logsql/tail?" + params.Encode()

	resp, err := c.httpClient.Get(ctx, fullPath)
	if err != nil {
		return &APIError{
			StatusCode: 0,
			Message:    fmt.Sprintf("tail request failed: %v", err),
			Query:      query,
		}
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    "tail request failed",
			Query:      query,
		}
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var raw map[string]interface{}
		if err := json.Unmarshal(line, &raw); err != nil {
			continue
		}

		entry := parseLogEntry(raw)
		if err := callback(entry); err != nil {
			return err
		}
	}

	return scanner.Err()
}

// TailWithLimit streams logs with entry limit
func (c *Client) TailWithLimit(ctx context.Context, query string, limit int) ([]LogEntry, error) {
	var entries []LogEntry
	count := 0

	// Create cancellable context
	tailCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := c.Tail(tailCtx, query, func(entry LogEntry) error {
		entries = append(entries, entry)
		count++
		if count >= limit {
			cancel() // Cancel when limit reached
			return nil
		}
		return nil
	})

	// If cancelled due to limit, not an error
	if err == context.Canceled && count >= limit {
		return entries, nil
	}

	return entries, err
}

// TailWithTimeout streams logs with timeout
func (c *Client) TailWithTimeout(ctx context.Context, query string, timeout time.Duration) ([]LogEntry, error) {
	tailCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var entries []LogEntry

	err := c.Tail(tailCtx, query, func(entry LogEntry) error {
		entries = append(entries, entry)
		return nil
	})

	// If cancelled due to timeout, not an error
	if err == context.DeadlineExceeded {
		return entries, nil
	}

	return entries, err
}

// parseLogEntry parses log entry
func parseLogEntry(raw map[string]interface{}) LogEntry {
	entry := LogEntry{
		Fields: make(map[string]interface{}),
	}

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

	for k, v := range raw {
		if k != "_time" && k != "_msg" && k != "_stream" {
			entry.Fields[k] = v
		}
	}

	return entry
}
