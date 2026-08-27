# design.md — shopify-mcp

## Class

**Tool-server** (not a relay). Does not implement WebSocket device pairing, edge catalogs, or hub correlation.

## Decisions

1. **Go** Streamable HTTP MCP (ported session pattern from `datumbridge-mcp-ws-hub`) instead of Python FastMCP — explicit platform exception for this service.
2. **GraphQL Admin API** only (`2025-01`), not REST Admin.
3. **Partner OAuth** lives in `datumbridge-mcp` vault (Office365-style). Tool-server only consumes `credentials_json`.
4. Writes require `confirm=true` or `dry_run=true` (Trello pattern).
5. Domain failures return MCP `isError=true`; protocol errors use JSON-RPC codes.

## Non-goals

Webhooks, Storefront API, Billing, Shopify Functions, `read_all_orders` Partner approval.
