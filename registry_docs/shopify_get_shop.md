# shopify_get_shop

Get shop identity, currency, and domain

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Audience

Curated store inquiry (ADR-0179). Anonymous / public DE turns may call this tool; Magento/Shopify/TikTok Shop still inject the workspace `wscred` login (ADR-0176).

| Who | Token |
|---|---|
| anonymous | none |

## Parameters

| Name | Required | Meaning |
|---|---|---|
| — | — | This call takes no model arguments. |

## Cases

### Typical call

Input:

```json
{}
```

Output:

```json
{
  "content": [
    {
      "type": "text",
      "text": "[UNTRUSTED_SHOPIFY_DATA] {\"shop\": {\"name\": \"Example Shop\"}}"
    }
  ],
  "isError": false
}
```
