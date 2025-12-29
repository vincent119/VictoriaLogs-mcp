package util

import "testing"

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"test", 4, "test"},
		{"", 10, ""},
	}

	for _, tt := range tests {
		result := TruncateString(tt.input, tt.maxLen)
		if result != tt.expected {
			t.Errorf("TruncateString(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
		}
	}
}

func TestTruncateString_Truncated(t *testing.T) {
	input := "hello world this is a long string"
	result := TruncateString(input, 5)

	// Should start with first 5 chars
	if result[:5] != "hello" {
		t.Errorf("Expected to start with 'hello', got %q", result[:5])
	}

	// Should contain truncation suffix
	if len(result) <= 5 {
		t.Error("Expected truncation suffix")
	}
}

func TestTruncateLines(t *testing.T) {
	input := "line1\nline2\nline3\nline4\nline5"

	result := TruncateLines(input, 3)

	// Result should contain first 3 lines
	if len(result) == 0 {
		t.Error("Expected non-empty result")
	}

	// Result should indicate truncation
	if result == input {
		t.Error("Expected truncated result")
	}
}

func TestTruncateLines_NoTruncation(t *testing.T) {
	input := "line1\nline2"
	result := TruncateLines(input, 5)

	if result != input {
		t.Errorf("Expected no truncation, got %q", result)
	}
}

func TestTruncateSlice(t *testing.T) {
	input := []string{"a", "b", "c", "d", "e"}

	result, truncated := TruncateSlice(input, 3)
	if len(result) != 3 {
		t.Errorf("Expected 3 elements, got %d", len(result))
	}
	if !truncated {
		t.Error("Expected truncated=true")
	}
	if result[0] != "a" || result[2] != "c" {
		t.Errorf("Unexpected slice content: %v", result)
	}
}

func TestTruncateSlice_NoTruncation(t *testing.T) {
	input := []string{"a", "b"}

	result, truncated := TruncateSlice(input, 5)
	if len(result) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(result))
	}
	if truncated {
		t.Error("Expected truncated=false")
	}
}

func TestTruncateMapSlice(t *testing.T) {
	input := []map[string]interface{}{
		{"key": "value1"},
		{"key": "value2"},
		{"key": "value3"},
	}

	result := TruncateMapSlice(input, 2)

	if !result.Truncated {
		t.Error("Expected truncated=true")
	}
	if result.ReturnedLen != 2 {
		t.Errorf("Expected 2 returned, got %d", result.ReturnedLen)
	}
	if result.OriginalLen != 3 {
		t.Errorf("Expected original 3, got %d", result.OriginalLen)
	}
}
