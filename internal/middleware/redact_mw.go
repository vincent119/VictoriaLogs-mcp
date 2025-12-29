package middleware

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/vincent119/victorialogs-mcp/internal/policy"
)

// RedactMiddleware Redact 中介層
type RedactMiddleware struct {
	redactor *policy.Redactor
}

// NewRedactMiddleware 建立 Redact 中介層
func NewRedactMiddleware(cfg policy.RedactConfig) *RedactMiddleware {
	return &RedactMiddleware{
		redactor: policy.NewRedactor(cfg),
	}
}

// Handler 回傳中介層處理函數
func (m *RedactMiddleware) Handler() ToolMiddleware {
	return func(next ToolHandler) ToolHandler {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			result, err := next(ctx, request)
			if err != nil {
				return result, err
			}

			// 對結果進行 redact 處理
			if result != nil {
				result = m.redactResult(result)
			}

			return result, nil
		}
	}
}

// redactResult 對結果進行 redact 處理
func (m *RedactMiddleware) redactResult(result *mcp.CallToolResult) *mcp.CallToolResult {
	if result == nil || len(result.Content) == 0 {
		return result
	}

	// 處理每個 content 項目
	for i, content := range result.Content {
		// 嘗試取得文字內容
		if textContent, ok := content.(mcp.TextContent); ok {
			textContent.Text = m.redactor.Apply(textContent.Text)
			result.Content[i] = textContent
		}
	}

	return result
}

// RedactString 對字串進行 redact 處理
func (m *RedactMiddleware) RedactString(s string) string {
	return m.redactor.Apply(s)
}

// RedactJSON 對 JSON 字串進行 redact 處理
func (m *RedactMiddleware) RedactJSON(data []byte) ([]byte, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		// 如果不是 JSON，直接處理字串
		return []byte(m.redactor.Apply(string(data))), nil
	}

	redacted := m.redactor.ApplyToMap(obj)
	return json.Marshal(redacted)
}
