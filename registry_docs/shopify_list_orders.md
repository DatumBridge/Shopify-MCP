# shopify_list_orders

List orders (default last 60 days without read_all_orders)

The gateway injects `credentials_json` from the connected account. Do not invent a token or paste a secret into the arguments.

## Parameters

| Name | Required | Meaning |
|---|---|---|
| `first` | no | string |
| `after` | no | string |
| `query` | no | string |

## Cases

### Typical call

Input:

```json
{
  "first": 25,
  "query": "quarterly report"
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
