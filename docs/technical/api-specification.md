# API specification — shopify-mcp

## REST

### `GET /health`

Returns `{"status":"ok","service":"shopify-mcp"}`.

### `POST /mcp`

JSON-RPC 2.0 Streamable HTTP subset.

| Method | Session | Notes |
|--------|---------|-------|
| `initialize` | creates `Mcp-Session-Id` | protocol `2024-11-05` |
| `notifications/initialized` | optional | HTTP 202 |
| `tools/list` | required | |
| `tools/call` | required | `{name, arguments}` |

### Error matrix

| Code | When |
|------|------|
| `-32700` | Invalid JSON |
| `-32600` | Invalid request |
| `-32601` | Unknown method/tool |
| `-32602` | Invalid params |
| `-32000` | Missing/invalid session |

Shopify domain failures: `result.isError=true` with text content.

## Tools

See `internal/tools/registry.go` for the authoritative list and `inputSchema`.
Common args: `credentials_json`, `credentials_path`. Writes: `confirm`, `dry_run`.
