// Package tools 提供獨立的 MCP Tool 實作檔案
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/vincent119/victorialogs-mcp/internal/victorialogs"
)

// TailHandler vlogs-tail Tool 處理器
type TailHandler struct {
	client *victorialogs.Client
}

// NewTailHandler 建立 Tail 處理器
func NewTailHandler(client *victorialogs.Client) *TailHandler {
	return &TailHandler{client: client}
}

// Handle 處理 vlogs-tail 請求
func (h *TailHandler) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("無效的請求參數"), nil
	}

	query, ok := args["query"].(string)
	if !ok || query == "" {
		return mcp.NewToolResultError("缺少必填參數: query"), nil
	}

	// 取得 limit（預設 100）
	limit := 100
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}
	if limit > 1000 {
		limit = 1000 // 限制最大值
	}

	// 取得 timeout（預設 5 秒）
	timeout := 5 * time.Second
	if t, ok := args["timeout"].(float64); ok {
		timeout = time.Duration(t) * time.Second
	}
	if timeout > 30*time.Second {
		timeout = 30 * time.Second // 限制最大超時
	}

	// 使用超時執行 Tail
	entries, err := h.client.TailWithTimeout(ctx, query, timeout)
	if err != nil && err != context.DeadlineExceeded {
		return mcp.NewToolResultError(fmt.Sprintf("Tail 失敗: %v", err)), nil
	}

	// 限制結果數量
	if len(entries) > limit {
		entries = entries[:limit]
	}

	// 格式化結果
	result := struct {
		Count   int                     `json:"count"`
		Entries []victorialogs.LogEntry `json:"entries"`
	}{
		Count:   len(entries),
		Entries: entries,
	}

	output, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(output)), nil
}
