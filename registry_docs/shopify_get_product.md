# shopify_get_product

Get a product by GraphQL ID or legacy numeric id

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Audience

Curated store inquiry (ADR-0179). Anonymous / public DE turns may call this tool; Magento/Shopify/TikTok Shop still inject the workspace `wscred` login (ADR-0176).

| Who | Token |
|---|---|
| anonymous | none |

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `id` | yes | string |

## Cases

### Typical call

Input:

```json
{
  "id": "example-id"
}
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

### Missing `id`

The tool rejects the call and does not guess the missing value.

Input:

```json
{}
```

Output:

```json
{
  "success": false,
  "error": {
    "error_code": "invalid_argument",
    "error_message": "id is required",
    "retryable": false
  }
}
```
