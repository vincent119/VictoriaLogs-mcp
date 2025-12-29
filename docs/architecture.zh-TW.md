# 系統架構設計

## 概述

VictoriaLogs MCP Server 是一個中介層應用程式，旨在連接 AI 助手（如 Claude）與 VictoriaLogs 日誌數據庫。它遵循 Model Context Protocol (MCP) 標準，將 VictoriaLogs 的查詢能力封裝為標準化的工具 (Tools)。

## 系統組件

```mermaid
graph TD
    Client[AI Client (Claude/Cline)] -->|JSON-RPC| MCPServer[MCP Server]
    MCPServer -->|Parse Request| Middleware[Middleware Layer]
    Middleware -->|Enforce Policy| VLClient[VictoriaLogs Client]
    VLClient -->|HTTP API| VLogs[VictoriaLogs DB]
```

### 1. MCP Server Layer
- **職責**：處理 JSON-RPC 通訊、工具註冊、請求路由。
- **實現**：基於 `mark3labs/mcp-go` 庫。
- **傳輸協議**：支援 Stdio (預設) 與 TCP。

### 2. Middleware Layer (中介層)
負責在請求到達核心邏輯前執行安全與合規檢查：
- **Rate Limit Middleware**：限制請求速率。
- **Circuit Breaker Middleware**：防止過載導致的級聯故障。
- **Audit Middleware**：記錄所有工具調用日誌。
- **Redact Middleware**：對回應結果進行敏感資料遮蔽。

### 3. Policy Layer (策略層)
定義具體的安全規則：
- **Allowlist**：限制可查詢的 Stream。
- **RateLimiter**：Token Bucket 算法實現。
- **Redactor**：Regex 基底的敏感資料過濾。

### 4. VictoriaLogs Client
- **職責**：封裝 VictoriaLogs HTTP API。
- **功能**：LogsQL 查詢構造、結果解析、錯誤處理。

## 目錄結構

- `cmd/vlmcp`: 應用程式入口。
- `internal/mcp`: MCP Server 相關實作 (Server, Tools)。
- `internal/policy`: 安全策略定義。
- `internal/middleware`: 中介層實作。
- `internal/victorialogs`: VictoriaLogs API 客戶端。
- `configs`: 設定檔與範例。

## 依賴管理

- **Viper**: 設定管理。
- **Zap (zlogger)**: 高效能日誌。
- **Prometheus**: 監控指標。
- **OpenTelemetry**: 分布式追蹤。
