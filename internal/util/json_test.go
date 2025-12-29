package util

import (
	"bytes"
	"testing"
)

func TestJSONEncode(t *testing.T) {
	input := map[string]interface{}{
		"name": "test",
		"count": 42,
	}

	result, err := JSONEncode(input)
	if err != nil {
		t.Fatalf("Failed to encode: %v", err)
	}

	if len(result) == 0 {
		t.Error("Expected non-empty result")
	}
}

func TestJSONDecode(t *testing.T) {
	input := bytes.NewReader([]byte(`{"name":"test","count":42}`))
	var result map[string]interface{}

	err := JSONDecode(input, &result)
	if err != nil {
		t.Fatalf("Failed to decode: %v", err)
	}

	if result["name"] != "test" {
		t.Errorf("Expected name=test, got %v", result["name"])
	}
}

func TestJSONToMap(t *testing.T) {
	input := []byte(`{"key":"value","number":123}`)
	result, err := JSONToMap(input)
	if err != nil {
		t.Fatalf("Failed to convert: %v", err)
	}

	if result["key"] != "value" {
		t.Errorf("Expected key=value, got %v", result["key"])
	}
}

func TestJSONEncodeIndent(t *testing.T) {
	input := map[string]string{"key": "value"}
	result, err := JSONEncodeIndent(input)
	if err != nil {
		t.Fatalf("Failed to encode: %v", err)
	}

	// Should contain indentation
	if !bytes.Contains(result, []byte("  ")) {
		t.Error("Expected indented output")
	}
}
