# shopify-mcp

Go MCP **tool-server** for Shopify Admin (GraphQL). Partner OAuth tokens are stored in the DatumBridge credential vault (`datumbridge-mcp`) and injected as `credentials_json` on tool execute.

## Endpoints

| Method | Path | Purpose |
|--------|------|---------|
| GET | `/health` | Liveness |
| POST | `/mcp` | Streamable HTTP MCP (`initialize`, `tools/list`, `tools/call`) |

Protocol: `2024-11-05`. Session header: `Mcp-Session-Id` (required after initialize).

## Auth

Tools accept `credentials_json` (or `credentials_path`):

```json
{
  "type": "oauth",
  "shop": "example.myshopify.com",
  "token": "<access_token>",
  "scopes": ["read_products", "write_products", "..."]
}
```

Partner app secrets stay in `datumbridge-mcp` env only — never in vault inject payloads.

Register as `mcpServer=shopify` (or `shopify-mcp`) so the registry vault injects credentials when args omit them.

## Tools

Read: shop, products, collections, orders, customers, inventory, discounts.  
Write: product/customer/order/inventory/discount/fulfillment mutations require `confirm=true` or `dry_run=true`.

## Run

```bash
go run ./cmd/api
# SHOPIFY_MCP_PORT=8010
```

```bash
go test ./...
```

See [design.md](design.md) and [docs/](docs/).
