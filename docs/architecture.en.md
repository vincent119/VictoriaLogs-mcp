# System Architecture Design

## Overview

The VictoriaLogs MCP Server is a middleware application designed to connect AI assistants (such as Claude) with the VictoriaLogs database. It follows the Model Context Protocol (MCP) standard, encapsulating VictoriaLogs query capabilities into standardized Tools.

## System Components

```mermaid
graph TD
    Client[AI Client (Claude/Cline)] -->|JSON-RPC| MCPServer[MCP Server]
    MCPServer -->|Parse Request| Middleware[Middleware Layer]
    Middleware -->|Enforce Policy| VLClient[VictoriaLogs Client]
    VLClient -->|HTTP API| VLogs[VictoriaLogs DB]
```

### 1. MCP Server Layer

- **Responsibilities**: Handles JSON-RPC communication, tool registration, and request routing.
- **Implementation**: Based on the `mark3labs/mcp-go` library.
- **Transport Protocols**: Supports Stdio (default) and TCP.

### 2. Middleware Layer

Responsible for executing security and compliance checks before requests reach core logic:

- **Rate Limit Middleware**: Limits request rates.
- **Circuit Breaker Middleware**: Prevents cascading failures due to overload.
- **Audit Middleware**: Logs all tool invocations.
- **Redact Middleware**: Masks sensitive data in response results.

### 3. Policy Layer

Defines specific security rules:

- **Allowlist**: Restricts queryable Streams.
- **RateLimiter**: Token Bucket algorithm implementation.
- **Redactor**: Regex-based sensitive data filtering.

### 4. VictoriaLogs Client

- **Responsibilities**: Encapsulates the VictoriaLogs HTTP API.
- **Functions**: LogsQL query construction, result parsing, error handling.

## Directory Structure

- `cmd/vlmcp`: Application entry point.
- `internal/mcp`: MCP Server related implementations (Server, Tools).
- `internal/policy`: Security policy definitions.
- `internal/middleware`: Middleware implementations.
- `internal/victorialogs`: VictoriaLogs API client.
- `configs`: Configuration files and examples.

## Dependencies

- **Viper**: Configuration management.
- **Zap (zlogger)**: High-performance logging.
- **Prometheus**: Monitoring metrics.
- **OpenTelemetry**: Distributed tracing.
