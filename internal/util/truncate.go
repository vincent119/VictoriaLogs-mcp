package util

import (
	"strings"
	"unicode/utf8"
)

// TruncateConfig 截斷設定
type TruncateConfig struct {
	MaxLength int    // 最大長度
	Suffix    string // 截斷後綴
}

// DefaultTruncateConfig 預設截斷設定
var DefaultTruncateConfig = TruncateConfig{
	MaxLength: 10000,
	Suffix:    "... [truncated]",
}

// TruncateString 截斷字串
func TruncateString(s string, maxLen int) string {
	if utf8.RuneCountInString(s) <= maxLen {
		return s
	}

	runes := []rune(s)
	return string(runes[:maxLen]) + DefaultTruncateConfig.Suffix
}

// TruncateStringWithConfig 使用設定截斷字串
func TruncateStringWithConfig(s string, cfg TruncateConfig) string {
	if utf8.RuneCountInString(s) <= cfg.MaxLength {
		return s
	}

	runes := []rune(s)
	suffixLen := utf8.RuneCountInString(cfg.Suffix)
	return string(runes[:cfg.MaxLength-suffixLen]) + cfg.Suffix
}

// TruncateLines 截斷行數
func TruncateLines(s string, maxLines int) string {
	lines := strings.Split(s, "\n")
	if len(lines) <= maxLines {
		return s
	}

	truncated := lines[:maxLines]
	truncated = append(truncated, "... ["+strings.TrimPrefix(DefaultTruncateConfig.Suffix, "... ")+" - remaining lines omitted]")
	return strings.Join(truncated, "\n")
}

// TruncateSlice 截斷 slice
func TruncateSlice[T any](slice []T, maxLen int) ([]T, bool) {
	if len(slice) <= maxLen {
		return slice, false
	}
	return slice[:maxLen], true
}

// TruncateResult 截斷結果結構
type TruncateResult struct {
	Data        interface{} `json:"data"`
	Truncated   bool        `json:"truncated"`
	OriginalLen int         `json:"original_len,omitempty"`
	ReturnedLen int         `json:"returned_len"`
}

// TruncateMapSlice 截斷 map slice 並回傳結果
func TruncateMapSlice(data []map[string]interface{}, maxLen int) TruncateResult {
	originalLen := len(data)
	truncated := false

	if originalLen > maxLen {
		data = data[:maxLen]
		truncated = true
	}

	return TruncateResult{
		Data:        data,
		Truncated:   truncated,
		OriginalLen: originalLen,
		ReturnedLen: len(data),
	}
}
