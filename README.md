# VictoriaLogs MCP Server

[![Go Version](https://img.shields.io/github/go-mod/go-version/vincent119/victorialogs-mcp)](go.mod)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

一個基於 [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) 構建的 VictoriaLogs 伺服器，讓 AI 助手（如 Claude、Cline）能夠直接查詢與分析 VictoriaLogs 中的日誌數據。

## ✨ 特色

- **LogsQL 支援**：直接執行 LogsQL 查詢，支援 pipe 語法。
- **健康檢查**：監控 VictoriaLogs 連線狀態。
- **Schema 探索**：自動探索可用的 Log Streams 與 Fields。
- **統計分析**：取得特定時間範圍內的日誌統計。
- **安全防護**：
  - Rate Limiting (速率限制)
  - Circuit Breaker (熔斷機制)
  - Sensitive Data Redaction (敏感資料遮蔽)
- **多種傳輸**：支援 Stdio (預設) 與 TCP 傳輸模式。

## 🚀 快速開始

### 安裝

```bash
# 下載最新發布版本
curl -LO https://github.com/vincent119/victorialogs-mcp/releases/latest/download/vlmcp-darwin-arm64
chmod +x vlmcp-darwin-arm64
mv vlmcp-darwin-arm64 vlmcp
```

### 設定

建立 `config.yaml`：

```yaml
victorialogs:
  url: "http://localhost:9428"
  auth:
    type: "none"

policy:
  rate_limit:
    enabled: true
    requests_per_minute: 60
```

### 與 Claude Desktop 整合

編輯 `~/Library/Application Support/Claude/claude_desktop_config.json`：

```json
{
  "mcpServers": {
    "victorialogs": {
      "command": "/path/to/vlmcp",
      "args": ["--config", "/path/to/config.yaml"]
    }
  }
}
```

## 🛠️ 可用工具 (Tools)

| 工具名稱 | 描述 |
|----------|------|
| `vlogs-query` | 執行 LogsQL 查詢 |
| `vlogs-stats` | 查詢日誌統計資料 (Hits) |
| `vlogs-schema` | 探索 Streams 與 Fields |
| `vlogs-health` | 檢查伺服器健康狀態 |

## 📚 文件

- [架構設計](docs/architecture.zh-TW.md)
- [API 與工具說明](docs/api-tools.zh-TW.md)
- [安全性設定](docs/security.zh-TW.md)
- [發布流程](docs/release-process.zh-TW.md)

## 📦 開發

```bash
# 建置
make build

# 測試
make test

# 執行 Lint
make lint
```

## 授權

MIT License
