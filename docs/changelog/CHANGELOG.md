# Changelog


## 2026-08-16

### Added

- Per-tool DatumBridge capabilities on MCP `tools/list` `_meta` (`capabilities` and `datumbridge.capabilities`). After rebuild and Tool Registry Setup / re-publish, business capabilities fill automatically.

## 2026-08-03

### Changed

- Vault Shopify bundles store shop + access token (+ scopes) only — Partner `client_secret` is not persisted or injected.
- Tool-server rejects non-`*.myshopify.com` shop hosts; CORS is allowlist-only; discount create requires amount or percentage.
