# API and Tools Description

This document details the tools provided by the VictoriaLogs MCP Server and how to use them.

## vlogs-query

Executes LogsQL queries and returns log entries.

### Parameters

| Parameter | Type | Required | Description | Example |
| :--- | :--- | :--- | :--- | :--- |
| `query` | string | Yes | LogsQL query syntax | `error` or `_stream:{app="web"}` |
| `limit` | number | No | Max entries to return (default 1000) | `100` |
| `start` | string | No | Start time (RFC3339 or relative time) | `5m`, `2024-01-01T00:00:00Z` |
| `end` | string | No | End time (default now) | `now` |

### Response Example

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

Queries log statistics (Hits).

### Parameters

| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `query` | string | No | LogsQL filter condition |
| `start` | string | Yes | Start time |
| `end` | string | No | End time |

## vlogs-schema

Explores available Schema information in VictoriaLogs.

### Parameters

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `type` | string | `streams` (list streams), `fields` (list fields), `values` (list field values) |
| `field` | string | Specify field name when type=values |
| `limit` | number | Max number to return |

## vlogs-health

Checks server connection status.

- **No Parameters Required**
- **Response**: `{"status": "healthy"}` or error message.
