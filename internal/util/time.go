package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 相對時間正規表示式
var relativeTimeRegex = regexp.MustCompile(`^(\d+)(s|m|h|d|w)$`)

// ParseTime 解析時間字串
// 支援格式:
// - RFC3339: 2024-01-01T00:00:00Z
// - Unix timestamp: 1704067200
// - 相對時間: 5m, 1h, 24h, 7d, 1w
func ParseTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, fmt.Errorf("time string is empty")
	}

	// 嘗試解析為相對時間
	if matches := relativeTimeRegex.FindStringSubmatch(s); len(matches) == 3 {
		return parseRelativeTime(matches[1], matches[2])
	}

	// 嘗試解析為 RFC3339
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}

	// 嘗試解析為 Unix timestamp
	if ts, err := strconv.ParseInt(s, 10, 64); err == nil {
		return time.Unix(ts, 0), nil
	}

	// 嘗試其他常見格式
	formats := []string{
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time format: %s", s)
}

// parseRelativeTime 解析相對時間
func parseRelativeTime(value, unit string) (time.Time, error) {
	n, err := strconv.Atoi(value)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time value: %s", value)
	}

	now := time.Now()
	var duration time.Duration

	switch unit {
	case "s":
		duration = time.Duration(n) * time.Second
	case "m":
		duration = time.Duration(n) * time.Minute
	case "h":
		duration = time.Duration(n) * time.Hour
	case "d":
		duration = time.Duration(n) * 24 * time.Hour
	case "w":
		duration = time.Duration(n) * 7 * 24 * time.Hour
	default:
		return time.Time{}, fmt.Errorf("invalid time unit: %s", unit)
	}

	return now.Add(-duration), nil
}

// FormatTime 格式化時間為 RFC3339
func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// FormatDuration 格式化時間間隔
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}

// ParseDuration 解析時間間隔字串
func ParseDuration(s string) (time.Duration, error) {
	// 先嘗試標準 time.Duration 格式
	if d, err := time.ParseDuration(s); err == nil {
		return d, nil
	}

	// 嘗試解析簡化格式（如 "1d", "1w"）
	s = strings.TrimSpace(s)
	if len(s) < 2 {
		return 0, fmt.Errorf("invalid duration: %s", s)
	}

	unit := s[len(s)-1:]
	value := s[:len(s)-1]

	n, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid time value: %s", value)
	}

	switch unit {
	case "d":
		return time.Duration(n) * 24 * time.Hour, nil
	case "w":
		return time.Duration(n) * 7 * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("invalid time unit: %s", unit)
	}
}

// TimeRangeWithinLimit 檢查時間範圍是否在限制內
func TimeRangeWithinLimit(start, end time.Time, limit time.Duration) bool {
	return end.Sub(start) <= limit
}
