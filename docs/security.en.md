# Security Configuration Guide

The VictoriaLogs MCP Server features built-in multi-layer protection mechanisms to ensure security in production environments.

## 1. Authentication and Authorization

Currently supports connection authentication for VictoriaLogs via `config.yaml` configuration:

- **Basic Auth**: Username/Password.
- **Bearer Token**: Token authentication.

```yaml
victorialogs:
  auth:
    type: "basic"
    username: "admin"
    password: "secure_password"
```

## 2. Rate Limiting

Prevents abuse or DDOS attacks by setting a maximum number of requests per minute.

```yaml
policy:
  rate_limit:
    enabled: true
    requests_per_minute: 60  # Max 60 requests per minute
```

## 3. Circuit Breaker

Automatically pauses requests to protect the system when the backend VictoriaLogs experiences persistent errors.

```yaml
policy:
  circuit_breaker:
    enabled: true
    error_threshold: 5  # Trigger circuit breaker after 5 consecutive errors
    timeout: "30s"      # Circuit breaker state maintained for 30 seconds
```

## 4. Redaction (Sensitive Data Masking)

Automatically detects and masks sensitive information in responses.

**Default Redaction Rules**:

- Email Addresses
- IP Addresses
- Credit Card Numbers
- API Keys / Tokens

This feature is enforced via Middleware and cannot be bypassed by the client.

## 5. Allowlist

Restricts AI to querying only specific Log Streams, ensuring data isolation.

```yaml
policy:
  allowlist:
    enabled: true
    streams:
      - "app/production/*"
      - "k8s/namespace/default"
```
