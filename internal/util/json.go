package util

import (
	"encoding/json"
	"fmt"
	"io"
)

// JSONDecode 解碼 JSON
func JSONDecode(r io.Reader, v interface{}) error {
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("JSON decode failed: %w", err)
	}
	return nil
}

// JSONEncode 編碼為 JSON
func JSONEncode(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("JSON encode failed: %w", err)
	}
	return data, nil
}

// JSONEncodeIndent 編碼為縮排 JSON
func JSONEncodeIndent(v interface{}) ([]byte, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("JSON encode failed: %w", err)
	}
	return data, nil
}

// JSONToMap 將 JSON 字串轉為 map
func JSONToMap(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("JSON decode failed: %w", err)
	}
	return result, nil
}

// MapToJSON 將 map 轉為 JSON 字串
func MapToJSON(m map[string]interface{}) ([]byte, error) {
	return JSONEncode(m)
}
