# shopify_adjust_inventory

Adjust available inventory at a location (confirm or dry_run)

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `inventory_item_id` | yes | string |
| `location_id` | yes | string |
| `delta` | yes | string |
| `reason` | no | string |
| `confirm` | no | Must be true to perform the write |
| `dry_run` | no | If true, return the planned request without calling Shopify |

## Cases

### Typical call

Input:

```json
{
  "inventory_item_id": "example-id",
  "location_id": "example-id",
  "delta": "example",
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

### Missing `inventory_item_id`

The tool rejects the call and does not guess the missing value.

Input:

```json
{
  "location_id": "example-id",
  "delta": "example",
  "dry_run": true
}
```

Output:

```json
{
  "success": false,
  "error": {
    "error_code": "invalid_argument",
    "error_message": "inventory_item_id is required",
    "retryable": false
  }
}
```

### Preview the write

Set `dry_run` to true. The tool returns the planned change and does not send it.

Input:

```json
{
  "inventory_item_id": "example-id",
  "location_id": "example-id",
  "delta": "example",
  "dry_run": true
}
```

### Confirmed write

Set `confirm` to true. Omit `dry_run`.

Input:

```json
{
  "inventory_item_id": "example-id",
  "location_id": "example-id",
  "delta": "example",
  "confirm": true
}
```

### Write without confirm or dry_run

Input:

```json
{
  "inventory_item_id": "example-id",
  "location_id": "example-id",
  "delta": "example"
}
```

Output:

```json
{
  "success": false,
  "error_message": "confirm=true required for write tools (or dry_run=true to preview)"
}
```
