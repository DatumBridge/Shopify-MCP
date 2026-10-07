# shopify_create_discount

Create a basic amount-off code discount (confirm or dry_run)

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `title` | yes | string |
| `code` | yes | string |
| `amount` | no | string |
| `percentage` | no | string |
| `starts_at` | no | string |
| `usage_limit` | no | string |
| `confirm` | no | Must be true to perform the write |
| `dry_run` | no | If true, return the planned request without calling Shopify |

## Cases

### Typical call

Input:

```json
{
  "title": "Status update",
  "code": "example",
  "dry_run": true
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

### Missing `title`

The tool rejects the call and does not guess the missing value.

Input:

```json
{
  "code": "example",
  "dry_run": true
}
```

Output:

```json
{
  "success": false,
  "error": {
    "error_code": "invalid_argument",
    "error_message": "title is required",
    "retryable": false
  }
}
```

### Preview the write

Set `dry_run` to true. The tool returns the planned change and does not send it.

Input:

```json
{
  "title": "Status update",
  "code": "example",
  "dry_run": true
}
```

### Confirmed write

Set `confirm` to true. Omit `dry_run`.

Input:

```json
{
  "title": "Status update",
  "code": "example",
  "confirm": true
}
```

### Write without confirm or dry_run

Input:

```json
{
  "title": "Status update",
  "code": "example"
}
```

Output:

```json
{
  "success": false,
  "error_message": "confirm=true required for write tools (or dry_run=true to preview)"
}
```
