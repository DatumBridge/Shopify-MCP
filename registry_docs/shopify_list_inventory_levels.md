# shopify_list_inventory_levels

List inventory levels for an inventory item

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `inventory_item_id` | yes | string |
| `first` | no | string |

## Cases

### Typical call

Input:

```json
{
  "inventory_item_id": "example-id",
  "first": 25
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
  "first": 25
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
