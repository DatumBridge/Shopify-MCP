# shopify_create_customer

Create a customer (confirm or dry_run)

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `email` | no | string |
| `first_name` | no | string |
| `last_name` | no | string |
| `phone` | no | string |
| `confirm` | no | Must be true to perform the write |
| `dry_run` | no | If true, return the planned request without calling Shopify |

## Cases

### Typical call

Input:

```json
{
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

### Preview the write

Set `dry_run` to true. The tool returns the planned change and does not send it.

Input:

```json
{
  "dry_run": true
}
```

### Confirmed write

Set `confirm` to true. Omit `dry_run`.

Input:

```json
{
  "confirm": true
}
```

### Write without confirm or dry_run

Input:

```json
{}
```

Output:

```json
{
  "success": false,
  "error_message": "confirm=true required for write tools (or dry_run=true to preview)"
}
```
