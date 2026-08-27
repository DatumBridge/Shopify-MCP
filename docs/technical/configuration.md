# Configuration — shopify-mcp

| Env | Default | Purpose |
|-----|---------|---------|
| `SHOPIFY_MCP_PORT` / `PORT` | `8010` | Listen port |
| `LOG_LEVEL` | `INFO` | DEBUG/INFO/WARN/ERROR |
| `SHOPIFY_MCP_ALLOWED_ORIGINS` | (empty) | Comma-separated browser Origins for CORS; empty disables CORS headers |
| `SHOPIFY_CREDENTIALS_DIR` | `/credentials` | Jail root for `credentials_path` reads |

Shopify shop + token come from tool args (`credentials_json`), not process env. Deploy behind the Tool Registry (cluster-internal); `/mcp` session is not identity.

Registry-side OAuth env (**datumbridge-integrations**, not MCP): `SHOPIFY_API_KEY`, `SHOPIFY_API_SECRET`, `STUDIO_PUBLIC_URL`, optional `CREDENTIAL_OAUTH_REDIRECT_URI_SHOPIFY`.
Allowlisted callback: `{STUDIO_PUBLIC_URL}/api/integrations/api/v1/credentials/oauth/shopify/callback`.
Partner **App URL** must use that same host (for local NodePort: `http://localhost:30080`). Do not set App URL to the shop (`*.myshopify.com`) or Shopify returns `redirect_uri and application url must have matching hosts`.

Install is blocked with “select a distribution method first” until the Partner app has a distribution. For Weaver Connect (one store / Plus org, not App Store) choose **Custom distribution**, add the shop `*.myshopify.com`, **Generate link**, then retry Studio Connect or open that install link. This choice cannot be changed later. See [Shopify: select a distribution method](https://shopify.dev/docs/apps/launch/distribution/select-distribution-method).
