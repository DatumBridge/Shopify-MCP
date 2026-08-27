# ADR-0001 — GraphQL Admin + vault OAuth for shopify-mcp

## Context

DatumBridge needs Shopify Admin automation via MCP. Platform SaaS MCP siblings are Python FastMCP; Go is used for the WS hub. Shopify Partner OAuth must be multi-user and vault-backed.

## Decision

1. Implement `shopify-mcp` as a Go tool-server with hand-rolled Streamable HTTP.
2. Use GraphQL Admin API exclusively.
3. Store OAuth tokens in `datumbridge-mcp` vault as `ShopifyTokenBundle` (includes shop domain).
4. Require confirm/dry_run for mutations.

## Alternatives Considered

1. Python FastMCP sibling — faster scaffold, but Go was required.
2. Custom Admin API token only — insufficient for Partner multi-store OAuth.
3. REST Admin — legacy; new public apps must use GraphQL.

## Consequences

- Broader vault provider surface; Studio Connect needs shop subdomain on start.
- No refresh_token failure gate (offline tokens are long-lived until revoked).

## Trade-offs

- More code than FastMCP; better alignment with existing Go hub HTTP patterns.

## Risks

- Shopify API version drift; pin `2025-01` and bump deliberately.
