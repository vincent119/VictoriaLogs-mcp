# API 與工具說明

本文件詳細說明 VictoriaLogs MCP Server 提供的工具及其使用方式。

## vlogs-query

執行 LogsQL 查詢並返回日誌條目。

### 參數

| 參數名 | 類型 | 必填 | 描述 | 範例 |
|--------|------|------|------|------|
| `query` | string | 是 | LogsQL 查詢語法 | `error` 或 `_stream:{app="web"}` |
| `limit` | number | 否 | 返回最大條數 (預設 1000) | `100` |
| `start` | string | 否 | 開始時間 (RFC3339 或相對時間) | `5m`, `2024-01-01T00:00:00Z` |
| `end` | string | 否 | 結束時間 (預設 now) | `now` |

### 回應範例

```json
[
  {
    "_msg": "connection failed",
    "_time": "2024-12-29T10:00:00Z",
    "_stream": "{app=\"payment\"}",
    "level": "error"
  }
]
```

## vlogs-stats

查詢日誌統計資料 (Hits)。

### 參數

| 參數名 | 類型 | 必填 | 描述 |
|--------|------|------|------|
| `query` | string | 否 | LogsQL 過濾條件 |
| `start` | string | 是 | 開始時間 |
| `end` | string | 否 | 結束時間 |

## vlogs-schema

探索 VictoriaLogs 中的可依據 Schema 資訊。

### 參數

| 參數名 | 類型 | 描述 |
|--------|------|------|
| `type` | string | `streams` (列出流), `fields` (列出欄位), `values` (列出欄位值) |
| `field` | string | 當 type=values 時指定欄位名 |
| `limit` | number | 返回最大數量 |

## vlogs-health

檢查伺服器連線狀態。

- **不需參數**
- **回應**：`{"status": "healthy"}` 或錯誤訊息。
