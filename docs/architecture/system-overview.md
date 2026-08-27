# System overview — shopify-mcp

## Role

Shopify Admin MCP tool-server for DatumBridge. Agents/Studio call tools via the Tool Registry, which injects per-user OAuth credentials from the vault. Each `tools/list` descriptor includes `_meta.capabilities` so Tool Registry business capabilities fill automatically on deploy/Setup.

## Components

| Component | Responsibility |
|-----------|----------------|
| `cmd/api` | HTTP bootstrap: `/health`, `POST /mcp` |
| `internal/mcp` | Session store + JSON-RPC Streamable HTTP |
| `internal/shopify` | Credentials parse + GraphQL Admin client |
| `internal/tools` | Tool descriptors and handlers |

## Data flow

```text
Studio/Agent → datumbridge-mcp (execute + vault inject)
  → shopify-mcp POST /mcp tools/call
  → Shopify GraphQL Admin API
```

## Dependencies

- Downstream: Shopify Admin GraphQL (`{shop}/admin/api/2025-01/graphql.json`)
- Upstream: Tool Registry / Studio (auth, vault)

## Risks

- Over-scoped OAuth if Partner app scopes are too broad
- Order history limited to ~60 days without `read_all_orders`
- Session store is in-memory (single replica)
