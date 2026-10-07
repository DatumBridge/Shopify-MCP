# shopify_create_fulfillment

Create a fulfillment for an order (confirm or dry_run)

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `order_id` | yes | string |
| `line_item_ids` | no | string |
| `tracking_number` | no | string |
| `tracking_company` | no | string |
| `notify_customer` | no | string |
| `confirm` | no | Must be true to perform the write |
| `dry_run` | no | If true, return the planned request without calling Shopify |

## Cases

### Typical call

Input:

```json
{
  "order_id": "example-id",
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

### Missing `order_id`

The tool rejects the call and does not guess the missing value.

Input:

```json
{
  "dry_run": true
}
```

Output:

```json
{
  "success": false,
  "error": {
    "error_code": "invalid_argument",
    "error_message": "order_id is required",
    "retryable": false
  }
}
```

### Preview the write

Set `dry_run` to true. The tool returns the planned change and does not send it.

Input:

```json
{
  "order_id": "example-id",
  "dry_run": true
}
```

### Confirmed write

Set `confirm` to true. Omit `dry_run`.

Input:

```json
{
  "order_id": "example-id",
  "confirm": true
}
```

### Write without confirm or dry_run

Input:

```json
{
  "order_id": "example-id"
}
```

Output:

```json
{
  "success": false,
  "error_message": "confirm=true required for write tools (or dry_run=true to preview)"
}
```
