# MCP Tool Specification — shopify-mcp

## Identity

| Field | Value |
|-------|-------|
| mcpServer | `shopify` (also matches `shopify-mcp`) |
| Transport | Streamable HTTP `POST /mcp` |
| Auth | Vault-injected `credentials_json` (Partner OAuth) |

## Intents

- Read/write Shopify Admin entities: products, collections, orders, customers, inventory, discounts.
- Preview mutations via `dry_run=true`; require `confirm=true` for side effects.

## Security / governance

- Secrets never logged.
- Deny writes without confirm/dry_run.
- Scopes requested at OAuth install (see datumbridge-mcp `defaultShopifyScopes`).
- Non-goals: webhooks, Storefront API, billing, policy/approval tools.

## Versioning

- Initial: 1.0.0
- Admin API pin: `2025-01`
