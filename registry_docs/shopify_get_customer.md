# shopify_get_customer

Get a customer by id

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Audience

Privileged store tool (ADR-0179). Requires a valid customer-session PAT on published DE turns (`allowedToolIds`). Mutate / PII / inventory write stays off `publicToolIds`.

| Who | Token |
|---|---|
| external | required |

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
