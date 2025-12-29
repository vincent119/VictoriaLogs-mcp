# VictoriaLogs MCP Server（Go）專案模板

## 目標：提供 MCP tools 給 AI/Agent 查詢 VictoriaLogs（Query / Tail / Stats / Explain…），預設走 stdio transport，不依賴 Gin

## 目錄結構（Template）

```bash
.
├── cmd/
│   └── vlmcp/
│       └── main.go
├── internal/
│   ├── app/
│   │   ├── app.go
│   │   └── lifecycle.go
│   ├── config/
│   │   ├── config.go
│   │   └── env.go
│   ├── logging/
│   │   └── logger.go
│   ├── mcp/
│   │   ├── server/
│   │   │   ├── server.go
│   │   │   ├── transport_stdio.go
│   │   │   ├── transport_tcp.go
│   │   │   └── errors.go
│   │   ├── schema/
│   │   │   ├── tools.go
│   │   │   └── types.go
│   │   └── tools/
│   │       ├── tool_query.go
│   │       ├── tool_tail.go
│   │       ├── tool_stats.go
│   │       ├── tool_explain.go
│   │       ├── tool_health.go
│   │       └── tool_schema.go
│   ├── middleware/            # MCP Tool 中介層
│   │   ├── ratelimit_mw.go
│   │   ├── redact_mw.go
│   │   ├── audit_mw.go
│   │   └── circuit_breaker_mw.go
│   ├── policy/
│   │   ├── policy.go
│   │   ├── allowlist.go
│   │   ├── redact_rules.go
│   │   ├── circuit_breaker.go
│   │   └── rate_limit.go
│   ├── victorialogs/
│   │   ├── client.go
│   │   ├── query.go
│   │   ├── tail.go
│   │   ├── stats.go
│   │   ├── schema.go
│   │   ├── models.go
│   │   └── errors.go
│   ├── observability/
│   │   ├── metrics.go
│   │   └── tracing.go
│   └── util/
│       ├── httpclient.go
│       ├── json.go
│       ├── time.go
│       └── truncate.go
├── pkg/
│   └── version/
│       └── version.go
├── configs/
│   ├── config.example.yaml
│   └── policy.example.yaml
├── docs/
│   ├── README.zh-TW.md
│   ├── README.en.md
│   ├── architecture.zh-TW.md
│   ├── architecture.en.md
│   ├── api-tools.zh-TW.md
│   ├── api-tools.en.md
│   ├── security.zh-TW.md
│   ├── security.en.md
│   ├── release-process.zh-TW.md
│   └── release-process.en.md
├── scripts/
│   ├── dev_run.sh
│   └── release_build.sh
├── testdata/
│   └── victorialogs_samples.json
├── .golangci.yml
├── .gitignore
├── go.mod
├── Makefile
└── README.md
```

## 檔案目錄清單

### 入口與生命週期

- [ ] cmd/vlmcp/main.go
- [ ] internal/app/app.go
- [ ] internal/app/lifecycle.go

### 設定與日誌

- [ ] internal/config/config.go
- [ ] internal/config/env.go
- [ ] internal/logging/logger.go

### MCP 協議與 Tool 定義

- [ ] internal/mcp/server/server.go
- [ ] internal/mcp/server/transport_stdio.go
- [ ] internal/mcp/server/transport_tcp.go
- [ ] internal/mcp/server/errors.go
- [ ] internal/mcp/schema/tools.go
- [ ] internal/mcp/schema/types.go

### MCP Tools（每個 Tool 一個檔案）

- [ ] internal/mcp/tools/tool_query.go
- [ ] internal/mcp/tools/tool_tail.go
- [ ] internal/mcp/tools/tool_stats.go
- [ ] internal/mcp/tools/tool_explain.go
- [ ] internal/mcp/tools/tool_health.go
- [ ] internal/mcp/tools/tool_schema.go

### 安全與治理（必備）

- [ ] internal/policy/policy.go
- [ ] internal/policy/allowlist.go
- [ ] internal/policy/rate_limit.go
- [ ] internal/policy/redact_rules.go
- [ ] internal/policy/circuit_breaker.go

### Middleware（MCP Tool 中介層）

- [ ] internal/middleware/ratelimit_mw.go
- [ ] internal/middleware/redact_mw.go
- [ ] internal/middleware/audit_mw.go
- [ ] internal/middleware/circuit_breaker_mw.go

### VictoriaLogs Client（對接 Query API）

- [ ] internal/victorialogs/client.go
- [ ] internal/victorialogs/query.go
- [ ] internal/victorialogs/tail.go
- [ ] internal/victorialogs/stats.go
- [ ] internal/victorialogs/schema.go
- [ ] internal/victorialogs/models.go
- [ ] internal/victorialogs/errors.go

### 可觀測性（建議）

- [ ] internal/observability/metrics.go
- [ ] internal/observability/tracing.go

### 通用工具

- [ ] internal/util/httpclient.go
- [ ] internal/util/json.go
- [ ] internal/util/time.go
- [ ] internal/util/truncate.go

### 版本與發布

- [ ] pkg/version/version.go
- [ ] scripts/dev_run.sh
- [ ] scripts/release_build.sh

### 設定範例

- [ ] configs/config.example.yaml
- [ ] configs/policy.example.yaml

### 文件（中英雙版）

- [ ] docs/README.zh-TW.md
- [ ] docs/README.en.md
- [ ] docs/architecture.zh-TW.md
- [ ] docs/architecture.en.md
- [ ] docs/api-tools.zh-TW.md
- [ ] docs/api-tools.en.md
- [ ] docs/security.zh-TW.md
- [ ] docs/security.en.md
  - 安全設計、Redact/Allowlist、審計與權限邊界
- [ ] docs/release-process.zh-TW.md
- [ ] docs/release-process.en.md
  - 版本策略、Tagging、Build/Release、回滾流程
- [ ] README.md

### 品質與自動化

- [ ] .golangci.yml
- [ ] Makefile
- [ ] .gitignore
- [ ] testdata/victorialogs_samples.json

---

## 🔧 變更細節（Change Details）

### 1. 核心與 MCP Server

| 檔案 | 描述 |
| ---- | ---- |
| `cmd/vlmcp/main.go` | 整合 Config、Logger、VictoriaLogs Client 並啟動 MCP Server |
| `internal/mcp/server/server.go` | 初始化 `mcp.Server`，註冊所有 Tools |
| `internal/mcp/server/transport_stdio.go` | 實作 Stdio transport（用於 Claude Desktop） |
| `internal/mcp/tools/*.go` | 各 Tool 的具體實作邏輯 |

### 2. Tools 參數設計

```go
// vlogs-query Tool
mcp.NewTool("vlogs-query",
    mcp.WithDescription("Execute LogsQL query against VictoriaLogs"),
    mcp.WithString("query", mcp.Required(), mcp.Description("LogsQL query string")),
    mcp.WithNumber("limit", mcp.Description("Maximum number of log entries to return")),
    mcp.WithString("start", mcp.Description("Start time (RFC3339 or relative like '5m')")),
    mcp.WithString("end", mcp.Description("End time (RFC3339 or relative)")),
)

// vlogs-stats Tool
mcp.NewTool("vlogs-stats",
    mcp.WithDescription("Get log statistics over time range"),
    mcp.WithString("query", mcp.Description("Optional filter query")),
    mcp.WithString("start", mcp.Required(), mcp.Description("Start time")),
    mcp.WithString("end", mcp.Description("End time (default: now)")),
)

// vlogs-schema Tool
mcp.NewTool("vlogs-schema",
    mcp.WithDescription("Explore available log streams and fields"),
    mcp.WithString("type", mcp.Enum("streams", "fields", "values"), mcp.Description("Type of schema info")),
    mcp.WithString("field", mcp.Description("Field name for values lookup")),
)
```

### 3. VictoriaLogs API 端點對應

| Tool | VictoriaLogs API |
| ---- | ---------------- |
| `vlogs-query` | `GET /select/logsql/query` |
| `vlogs-tail` | `GET /select/logsql/tail` |
| `vlogs-stats` | `GET /select/logsql/hits` |
| `vlogs-schema` (streams) | `GET /select/logsql/streams` |
| `vlogs-schema` (fields) | `GET /select/logsql/field_names` |
| `vlogs-explain` | `GET /select/logsql/query`（附加 explain 類參數，依 VictoriaLogs 版本為準） |
| `vlogs-health` | `GET /health` |

### 4. 設定檔結構（config.example.yaml）

```yaml
server:
  name: "victorialogs-mcp"
  version: "1.0.0"
  transport: "stdio"        # stdio | tcp
  tcp_addr: ":9090"         # 僅當 transport=tcp 時使用

victorialogs:
  url: "http://localhost:9428"
  auth:
    type: "none"            # none | basic | bearer
    username: ""
    password: ""
    token: ""
  timeout: "30s"            # HTTP 請求超時
  query_timeout: "60s"      # 單次查詢最大執行時間
  max_results: 5000

policy:
  rate_limit:
    enabled: true
    requests_per_minute: 60
  allowlist:
    enabled: false
    streams: []             # 允許查詢的 stream patterns
  circuit_breaker:
    enabled: true
    error_threshold: 5      # 連續 N 次錯誤觸發熔斷
    timeout: "30s"          # 熔斷持續時間

logging:
  level: "info"             # debug | info | warn | error
  format: "json"            # json | text
```

### 5. Graceful Shutdown 考量（lifecycle.go）

- 等待進行中的查詢完成（Context Timeout）
- 關閉 VictoriaLogs Client HTTP 連線池
- Flush 審計日誌 buffer
- 釋放 Rate Limiter 資源

---

## ✅ 驗證計畫（Verification Plan）

### 自動化測試

| 測試類型 | 涵蓋範圍 |
| -------- | -------- |
| Unit Tests | `internal/victorialogs` Client 邏輯（使用 httptest Mock） |
| Tool Tests | MCP Tool Handler 參數解析與錯誤處理 |
| Policy Tests | Rate Limit、Redact、Allowlist 規則驗證 |
| Integration | 端到端流程（需要實際 VictoriaLogs 實例或 Mock Server） |

**執行指令**：

```bash
// turbo
make test

// turbo
make lint
```

### 手動驗證

1. **建置**：

   ```bash
   // turbo
   make build
   ```

2. **本地測試**：

   ```bash
   ./bin/vlmcp --config configs/config.example.yaml
   ```

3. **Claude Desktop 整合測試**：
   - 編輯 `~/Library/Application Support/Claude/claude_desktop_config.json`：

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

   - 對 Claude 下達指令：「查詢過去 5 分鐘含有 'error' 的日誌」
   - 確認 AI 能正確呼叫 `vlogs-query` 並解讀回傳的日誌

4. **MCP Inspector 測試**：

   ```bash
   npx @modelcontextprotocol/inspector ./bin/vlmcp
   ```

---

## ⚠️ 使用者審查項目（User Review Required）

> 以下項目需由專案負責人或系統擁有者確認，
> 目的在於避免 MCP Server 在生產環境中造成 **資料外洩、資源濫用、或 AI 誤用風險**。

---

### 1️⃣ 依賴管理（Dependency Management）

- **MCP SDK**
  - 使用社群標準 `github.com/mark3labs/mcp-go`
  - 鎖定版本以避免協議升級造成相容性問題

- **Config**
  - 使用 `spf13/viper`
  - 允許 runtime reload 設定
  - 必須限制敏感設定來源（env / file / secret manager）

- **Logging**
  - 使用標準化 Logger `github.com/vincent119/zlogger`
  - 必須避免記錄完整 LogsQL query
  - 必須避免記錄原始 log payload（raw logs）

---

### 2️⃣ MCP Tool 行為與權限審查（Tool Governance）

- **可啟用的 MCP Tools 是否受限**
  - 允許查詢日誌：`vlogs-query`
  - 允許統計分析：`vlogs-stats`
  - 允許查詢結構（Streams / Fields）：`vlogs-schema`
  - 允許即時串流：`vlogs-tail`
  - 允許查詢執行計畫：`vlogs-explain`
  - 允許健康檢查：`vlogs-health`

- **生產環境高風險工具管控**
  - 是否允許即時串流工具：`vlogs-tail`
  - 是否允許查詢執行計畫工具：`vlogs-explain`

- **查詢與回傳限制**
  - 單次 Tool 呼叫最大回傳筆數：`5000`
  - 單次查詢允許的最長時間範圍：`6 小時以內`

---

### 3️⃣ 查詢風險控管（Query Safety）

- **LogsQL 查詢複雜度限制**
  - 是否禁止無任何 filter 的全量掃描（Full Scan）
  - 是否強制查詢必須包含時間條件（start / end 或相對時間）

- **查詢資源限制設定**
  - 最大查詢時間範圍（例如：`6 小時`）
  - 單次查詢最大結果筆數（例如：`5,000 筆`）
  - 單一使用者 / Agent 的查詢頻率上限（Rate Limit）

- **高成本查詢額外管控**
  - 是否對 `vlogs-explain` 套用更嚴格的限制
  - 是否對 `vlogs-stats` 套用更嚴格的限制
  - 是否限制 Explain / Stats 僅能在非尖峰時段使用

---

### 4️⃣ 資料存取控制與敏感資訊防護（Access Control & Redaction）

- **Allowlist（白名單）機制**
  - 是否啟用 Allowlist 限制可查詢的 log streams / namespaces
  - 是否明確禁止查詢不相關服務或系統層級日誌
  - 是否防止 AI 掃描敏感系統（如 auth、infra、security）相關日誌

- **Redact（遮罩）規則**
  - 是否對以下敏感資訊進行遮罩或移除：
    - IP 位址（IPv4 / IPv6）
    - Authorization Header、Bearer Token、Cookie
    - Email、帳號、使用者識別碼（User ID）
    - API Key、Session ID、Trace ID（視情境）

- **回傳內容安全控管**
  - 是否避免直接回傳原始 Request / Response Headers
  - 是否避免回傳內部網路資訊：
    - 內部 IP（Private IP）
    - Pod IP、Node IP
    - Service Mesh / Overlay Network 位址

---

### 5️⃣ MCP Transport 與網路暴露風險（Transport & Exposure Control）

- **Transport 模式管控**
  - 是否預設僅啟用 `stdio` transport（本機、單一使用者）
  - 是否明確區分 `stdio` 與 `tcp` transport 的使用場景
  - 是否禁止在未授權情境下啟用 `tcp` transport

- **TCP 模式風險評估**
  - 是否限制 TCP 監聽位址（僅允許 `127.0.0.1` 或內網）
  - 是否禁止直接暴露於公網（0.0.0.0）
  - 是否搭配防火牆或安全群組限制來源 IP

- **身份與存取控制**
  - 是否為 TCP transport 加入額外驗證機制（Token / mTLS）
  - 是否限制可連線的 Client / Agent 清單
  - 是否記錄每個連線的來源與存取行為

- **AI Client 使用風險**
  - 是否限制僅允許受信任的 AI Client（如 Claude Desktop）
  - 是否避免多個 AI Agent 共用同一 MCP Server 實例
  - 是否防止未經授權的 Tool 呼叫與重放請求（Replay）

- **錯誤與異常暴露**
  - 是否避免將內部錯誤訊息直接回傳給 AI
  - 是否統一錯誤格式，避免洩漏內部結構與實作細節

---

### 6️⃣ 審計、稽核與 AI 行為追蹤（Audit & Observability）

- **MCP Tool 呼叫審計**
  - 是否記錄每一次 MCP Tool 呼叫行為
  - 是否記錄：
    - Tool 名稱
    - 呼叫時間
    - 呼叫來源（AI Client / Agent）
    - 查詢參數摘要（避免記錄完整 LogsQL）
  - 是否避免在審計紀錄中保存原始 log payload
  - 審計與行為紀錄僅保留必要欄位，遵循最小化與最短保留原則

- **AI 行為可追蹤性**
  - 是否可追溯：
    - 哪個 AI / Agent 觸發查詢
    - 查詢了哪些 streams / namespaces
    - 是否觸發限制或被拒絕
  - 是否為每次請求產生唯一 Trace ID

- **異常與濫用偵測**
  - 是否偵測異常行為：
    - 短時間大量查詢
    - 多次觸發限制（rate limit / time range）
    - 嘗試繞過 Allowlist / Redact 規則
  - 是否對異常行為進行告警（Alert）

- **可觀測性指標（Metrics）**
  - MCP Tool 呼叫次數（依 Tool 類型）
  - 查詢成功 / 失敗比率
  - 查詢平均延遲與 P95 / P99
  - 被拒絕的查詢數（Policy / Rate Limit）
  - Redact 命中次數

- **日誌與指標儲存策略**
  - 是否將 MCP Server 審計日誌送往集中式 Log 系統
  - 是否將 Metrics 暴露給 Prometheus
  - 是否設定合理的審計資料保留期限（Retention）

- **合規與責任界線**
  - 是否清楚定義：
    - AI 只能「輔助分析」
    - 不具備直接資料存取或決策權限
  - 是否有文件說明 MCP Server 的使用責任歸屬

---

### 7️⃣ 生產環境啟用前確認（Go / No-Go Checklist）

> 以下項目必須在 MCP Server 進入生產環境前逐項確認完成，任一項未通過即視為 **No-Go**。

- [ ] **Rate Limit 已啟用**
  - 已限制單一使用者 / Agent 的查詢頻率
  - 已驗證超額請求會被正確拒絕

- [ ] **Redact 規則已啟用**
  - 已遮罩 IP、Token、Authorization、Cookie
  - 已確認不回傳原始敏感欄位

- [ ] **查詢資源限制已設定**
  - 已限制最大查詢時間範圍
  - 已限制單次回傳最大筆數

- [ ] **MCP Tools 白名單已確認**
  - 已明確啟用允許的 Tools
  - 已關閉或限制高風險 Tools（如 `vlogs-tail`、`vlogs-explain`）

- [ ] **異常行為與錯誤處理已測試**
  - 已測試無 filter、超大時間範圍等異常查詢
  - 已確認錯誤回傳不洩漏內部實作細節

- [ ] **人工安全審查已完成**
  - 已由系統擁有者或資安角色完成審查
  - 已確認 MCP Server 僅作為分析輔助，不具直接資料存取權限

- [ ] **Transport 模式已確認**
  - Production 環境僅允許 `stdio` 或 `tcp` 限制於 `127.0.0.1` / 內網
  - 已確認未暴露 `0.0.0.0`
  - 已驗證 TCP 模式搭配防火牆或安全群組限制來源 IP
  - 驗證方式：`ss -lntp | grep <port>` 或 `netstat -an | grep LISTEN`

---

### 🧠 補充說明

> MCP Server 本質上是 **AI 的能力放大器**。
> 本審查清單的目標在於確保：
>
> - AI 只能看到它應該看到的資料
> - AI 不能執行人類管理者不會允許的操作
>
> 本文件適用於 VictoriaLogs MCP Server v1.x
> MCP Tool 行為變更需同步更新本文件與 ADR
