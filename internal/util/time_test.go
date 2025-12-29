package util

import (
	"testing"
	"time"
)

func TestParseTime_RFC3339(t *testing.T) {
	input := "2024-01-15T10:30:00Z"
	result, err := ParseTime(input)
	if err != nil {
		t.Fatalf("Failed to parse RFC3339: %v", err)
	}

	expected := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestParseTime_UnixTimestamp(t *testing.T) {
	input := "1704067200"
	result, err := ParseTime(input)
	if err != nil {
		t.Fatalf("Failed to parse Unix timestamp: %v", err)
	}

	expected := time.Unix(1704067200, 0)
	if !result.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestParseTime_Relative(t *testing.T) {
	tests := []struct {
		input    string
		minAgo   time.Duration
		maxAgo   time.Duration
	}{
		{"5m", 4*time.Minute + 50*time.Second, 5*time.Minute + 10*time.Second},
		{"1h", 59*time.Minute + 50*time.Second, 60*time.Minute + 10*time.Second},
		{"24h", 23*time.Hour + 59*time.Minute, 24*time.Hour + 1*time.Minute},
	}

	for _, tt := range tests {
		result, err := ParseTime(tt.input)
		if err != nil {
			t.Errorf("Failed to parse %s: %v", tt.input, err)
			continue
		}

		ago := time.Since(result)
		if ago < tt.minAgo || ago > tt.maxAgo {
			t.Errorf("For %s, expected ~%v ago, got %v ago", tt.input, tt.minAgo, ago)
		}
	}
}

func TestParseTime_Empty(t *testing.T) {
	_, err := ParseTime("")
	if err == nil {
		t.Error("Expected error for empty string")
	}
}

func TestParseTime_Invalid(t *testing.T) {
	_, err := ParseTime("invalid-time")
	if err == nil {
		t.Error("Expected error for invalid time")
	}
}

func TestFormatTime(t *testing.T) {
	input := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	result := FormatTime(input)
	expected := "2024-01-15T10:30:00Z"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		input    time.Duration
		expected string
	}{
		{30 * time.Second, "30s"},
		{5 * time.Minute, "5m"},
		{2 * time.Hour, "2h"},
		{48 * time.Hour, "2d"},
	}

	for _, tt := range tests {
		result := FormatDuration(tt.input)
		if result != tt.expected {
			t.Errorf("For %v, expected %s, got %s", tt.input, tt.expected, result)
		}
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
	}{
		{"1d", 24 * time.Hour},
		{"1w", 7 * 24 * time.Hour},
		{"1h", time.Hour},
		{"30m", 30 * time.Minute},
	}

	for _, tt := range tests {
		result, err := ParseDuration(tt.input)
		if err != nil {
			t.Errorf("Failed to parse %s: %v", tt.input, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("For %s, expected %v, got %v", tt.input, tt.expected, result)
		}
	}
}

func TestTimeRangeWithinLimit(t *testing.T) {
	now := time.Now()

	// Within limit
	start := now.Add(-1 * time.Hour)
	if !TimeRangeWithinLimit(start, now, 2*time.Hour) {
		t.Error("1 hour range should be within 2 hour limit")
	}

	// Exceeds limit
	start = now.Add(-3 * time.Hour)
	if TimeRangeWithinLimit(start, now, 2*time.Hour) {
		t.Error("3 hour range should exceed 2 hour limit")
	}
}
