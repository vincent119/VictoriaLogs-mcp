# 安全性設定指南

VictoriaLogs MCP Server 內建多層防護機制，確保在生產環境中的安全性。

## 1. 認證與授權

目前支援通過 `config.yaml` 配置 VictoriaLogs 的連線認證：
- **Basic Auth**：使用者名稱/密碼。
- **Bearer Token**：Token 認證。

```yaml
victorialogs:
  auth:
    type: "basic"
    username: "admin"
    password: "secure_password"
```

## 2. Rate Limiting (速率限制)

防止濫用或 DDOS 攻擊，可設定每分鐘最大請求數。

```yaml
policy:
  rate_limit:
    enabled: true
    requests_per_minute: 60  # 每分鐘最多 60 次請求
```

## 3. Circuit Breaker (熔斷機制)

當後端 VictoriaLogs 出現持續錯誤時，自動暫停請求以保護系統。

```yaml
policy:
  circuit_breaker:
    enabled: true
    error_threshold: 5  # 連續 5 次錯誤觸發熔斷
    timeout: "30s"      # 熔斷狀態維持 30 秒
```

## 4. Redaction (敏感資料遮蔽)

自動偵測並遮蔽回應中的敏感資訊。

**預設遮蔽規則**：
- Email 地址
- IP 地址
- Credit Card 號碼
- API Keys / Tokens

此功能透過 Middleware 強制執行，無法被客戶端繞過。

## 5. Allowlist (白名單)

限制 AI 只能查詢特定的 Log Stream，確保數據隔離。

```yaml
policy:
  allowlist:
    enabled: true
    streams:
      - "app/production/*"
      - "k8s/namespace/default"
```
