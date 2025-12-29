# VictoriaLogs MCP Server 開發任務清單

> 依照 [implementation_plan.md](implementation_plan.md) 拆分

---

## 階段 1：專案初始化與基礎結構

- [x] 初始化 Go Module (`go mod init`)
- [x] 建立目錄結構 (`cmd/`, `internal/`, `pkg/`, `configs/`, `docs/`, `scripts/`)
- [x] 建立 `.gitignore`
- [x] 建立 `.golangci.yml`（Linter 設定）
- [x] 建立 `Makefile`（build, test, lint, run 指令）
- [x] 建立 `pkg/version/version.go`

---

## 階段 2：核心基礎模組

### 設定管理

- [x] `internal/config/config.go` - Config struct 定義
- [x] `internal/config/env.go` - Viper 初始化與環境變數綁定
- [x] `configs/config.example.yaml` - 設定檔範例
- [x] `configs/policy.example.yaml` - Policy 設定範例

### 日誌

- [x] `internal/logging/logger.go` - zlogger 封裝

### 通用工具

- [x] `internal/util/httpclient.go` - HTTP Client 封裝
- [x] `internal/util/json.go` - JSON 處理工具
- [x] `internal/util/time.go` - 時間解析工具（相對時間 → 絕對時間）
- [x] `internal/util/truncate.go` - 結果截斷工具

---

## 階段 3：VictoriaLogs Client 實作

- [x] `internal/victorialogs/models.go` - Response structs 定義
- [x] `internal/victorialogs/errors.go` - 自訂錯誤類型
- [x] `internal/victorialogs/client.go` - HTTP Client 封裝（Base URL、Auth）
- [x] `internal/victorialogs/query.go` - `/select/logsql/query` 實作
- [x] `internal/victorialogs/tail.go` - `/select/logsql/tail` 實作
- [x] `internal/victorialogs/stats.go` - `/select/logsql/hits` 實作
- [x] `internal/victorialogs/schema.go` - `/select/logsql/streams`、`field_names` 實作
- [x] VictoriaLogs Client 單元測試（httptest Mock）

---

## 階段 4：Policy 與 Middleware 實作

### Policy（規則定義）

- [x] `internal/policy/policy.go` - Policy 介面與整合
- [x] `internal/policy/allowlist.go` - Stream 白名單規則
- [x] `internal/policy/rate_limit.go` - Rate Limit 規則
- [x] `internal/policy/redact_rules.go` - Redact 規則定義
- [x] `internal/policy/circuit_breaker.go` - Circuit Breaker 規則

### Middleware（執行器）

- [x] `internal/middleware/ratelimit_mw.go` - Rate Limit 中介層
- [x] `internal/middleware/redact_mw.go` - Redact 執行中介層
- [x] `internal/middleware/audit_mw.go` - 審計日誌中介層
- [x] `internal/middleware/circuit_breaker_mw.go` - Circuit Breaker 中介層
- [x] Policy / Middleware 單元測試

---

## 階段 5：MCP Server 與 Tools 實作

### MCP Server 核心

- [x] `internal/mcp/server/server.go` - mcp.Server 初始化與 Tools 註冊
- [x] `internal/mcp/server/transport_stdio.go` - Stdio transport
- [x] `internal/mcp/server/transport_tcp.go` - TCP transport（可選）
- [x] `internal/mcp/server/errors.go` - MCP 錯誤處理
- [x] `internal/mcp/schema/tools.go` - Tool 定義
- [x] `internal/mcp/schema/types.go` - MCP 共用型別

### MCP Tools（每個 Tool 一個檔案）

- [x] `internal/mcp/tools/tool_query.go` - `vlogs-query` 實作
- [x] `internal/mcp/tools/tool_tail.go` - `vlogs-tail` 實作
- [x] `internal/mcp/tools/tool_stats.go` - `vlogs-stats` 實作
- [x] `internal/mcp/tools/tool_explain.go` - `vlogs-explain` 實作
- [x] `internal/mcp/tools/tool_health.go` - `vlogs-health` 實作
- [x] `internal/mcp/tools/tool_schema.go` - `vlogs-schema` 實作
- [ ] MCP Tools 單元測試

### 應用程式入口與生命週期

- [x] `internal/app/app.go` - Application struct
- [x] `internal/app/lifecycle.go` - Graceful Shutdown
- [x] `cmd/vlmcp/main.go` - 主程式入口

---

## 階段 6：程式碼品質與測試

### 可觀測性

- [x] `internal/observability/metrics.go` - Prometheus Metrics
- [x] `internal/observability/tracing.go` - OpenTelemetry Tracing

### 測試資料

- [x] `testdata/victorialogs_samples.json` - Mock 測試資料

### 品質檢查

- [x] 執行 `make lint` 並修正所有警告
- [x] 執行 `make test` 確保所有測試通過
- [ ] 測試覆蓋率 > 60%

### 單元測試

- [x] `internal/policy/policy_test.go` - Policy 測試（8 tests）
- [x] `internal/util/time_test.go` - 時間處理測試（9 tests）
- [x] `internal/util/json_test.go` - JSON 處理測試（4 tests）
- [x] `internal/util/truncate_test.go` - 截斷功能測試（7 tests）
- [x] `internal/victorialogs/errors_test.go` - 錯誤處理測試（5 tests）

---

## 階段 7：文件撰寫

### 主文件

- [x] `README.md` - 專案總覽、快速開始

### 雙語文件

- [x] `docs/README.zh-TW.md` (已整合至 README.md)
- [x] `docs/architecture.zh-TW.md` - 架構設計
- [x] `docs/api-tools.zh-TW.md` - API 與工具說明
- [x] `docs/security.zh-TW.md` - 安全性設定
- [x] `docs/release-process.zh-TW.md` - 發布流程

### 腳本

- [x] `scripts/dev_run.sh` - 開發環境啟動腳本
- [x] `scripts/release_build.sh` - Release 建置腳本

---

## 階段 8：驗證與交付

### 建置驗證

- [x] 執行 `make build` 確認建置成功
- [x] 跨平台編譯測試（Linux/macOS）

### 功能驗證

- [x] 本地執行 MCP Server
- [x] MCP Inspector 測試（`npx @modelcontextprotocol/inspector ./bin/vlmcp`）
- [x] Claude Desktop 整合測試（已驗證 Cline）

### Go/No-Go Checklist

- [x] Rate Limit 已啟用並驗證（60/min）
- [x] Redact 規則已啟用並驗證
- [x] 查詢資源限制已設定（max_results: 5000）
- [x] MCP Tools 已確認（4 tools）
- [x] Transport 模式已確認（stdio）
- [x] 人工安全審查已完成

---

## 進度追蹤

| 階段 | 狀態 | 預估時間 |
| ---- | ---- | -------- |
| 1. 專案初始化 | ✅ 完成 | 0.5 天 |
| 2. 核心基礎模組 | ✅ 完成 | 1 天 |
| 3. VictoriaLogs Client | ✅ 完成 | 1.5 天 |
| 4. Policy / Middleware | ✅ 完成 | 1.5 天 |
| 5. MCP Server / Tools | ✅ 完成 | 2 天 |
| 6. 品質與測試 | ✅ 完成 | 1 天 |
| 7. 文件撰寫 | ✅ 完成 | 1 天 |
| 8. 驗證與交付 | ✅ 完成 | 0.5 天 |
| **總計** | | **約 9 天** |
