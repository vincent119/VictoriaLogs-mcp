# VictoriaLogs MCP Server

[![Go Version](https://img.shields.io/github/go-mod/go-version/vincent119/victorialogs-mcp)](../go.mod)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](../LICENSE)

A VictoriaLogs server built on the [Model Context Protocol (MCP)](https://modelcontextprotocol.io/), enabling AI assistants (such as Claude, Cline) to directly query and analyze log data in VictoriaLogs.

## ✨ Features

- **LogsQL Support**: Execute LogsQL queries directly, supporting pipe syntax.
- **Health Check**: Monitor VictoriaLogs connection status.
- **Schema Discovery**: Automatically discover available Log Streams and Fields.
- **Statistical Analysis**: Retrieve log statistics for specific time ranges.
- **Security**:
  - Rate Limiting
  - Circuit Breaker
  - Sensitive Data Redaction
- **Multiple Transports**: Supports Stdio (default) and TCP transport modes.

## 🚀 Quick Start

### Installation

```bash
# Download the latest release
curl -LO https://github.com/vincent119/victorialogs-mcp/releases/latest/download/vlmcp-darwin-arm64
chmod +x vlmcp-darwin-arm64
mv vlmcp-darwin-arm64 vlmcp
```

### Configuration

Create `config.yaml`:

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

### Integration with Claude Desktop

Edit `~/Library/Application Support/Claude/claude_desktop_config.json`:

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

## 🛠️ Available Tools

| Tool Name | Description |
| :--- | :--- |
| `vlogs-query` | Execute LogsQL queries |
| `vlogs-stats` | Query log statistics (Hits) |
| `vlogs-schema` | Explore Streams and Fields |
| `vlogs-health` | Check server health status |

## 📚 Documentation

- [Architecture Design](architecture.en.md)
- [API and Tools Description](api-tools.en.md)
- [Security Configuration](security.en.md)
- [Release Process](release-process.en.md)

## 📦 Development

```bash
# Build
make build

# Test
make test

# Run Lint
make lint
```

## License

MIT License
