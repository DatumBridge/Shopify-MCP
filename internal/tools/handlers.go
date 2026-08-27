package tools

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/datumbridge/shopify-mcp/internal/mcp"
	"github.com/datumbridge/shopify-mcp/internal/shopify"
)

func handleGetShop(raw json.RawMessage) map[string]interface{} {
	c, _, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	ctx, cancel := ctx()
	defer cancel()
	const q = `query { shop { id name email myshopifyDomain primaryDomain { url host } currencyCode timezoneAbbreviation } }`
	data, err := c.GraphQL(ctx, q, nil)
	return gqlOK(data, err)
}

func handleListProducts(raw json.RawMessage) map[string]interface{} {
	c, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	ctx, cancel := ctx()
	defer cancel()
	vars := pageVars(m)
	const q = `query($first: Int!, $after: String, $query: String) {
	  products(first: $first, after: $after, query: $query) {
	    pageInfo { hasNextPage endCursor }
	    edges { cursor node { id title handle status vendor productType updatedAt
	      variants(first: 10) { edges { node { id title sku inventoryQuantity } } }
	    } }
	  }
	}`
	data, err := c.GraphQL(ctx, q, vars)
	return gqlOK(data, err)
}

func handleGetProduct(raw json.RawMessage) map[string]interface{} {
	c, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	id := productGID(strArg(m, "id"))
	if id == "gid://shopify/Product/" {
		return mcp.ToolResultError("id required")
	}
	ctx, cancel := ctx()
	defer cancel()
	const q = `query($id: ID!) {
	  product(id: $id) {
	    id title handle descriptionHtml status vendor productType tags updatedAt
	    variants(first: 50) { edges { node { id title sku price inventoryQuantity inventoryItem { id } } } }
	  }
	}`
	data, err := c.GraphQL(ctx, q, map[string]interface{}{"id": id})
	return gqlOK(data, err)
}

func handleCreateProduct(raw json.RawMessage) map[string]interface{} {
	c, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	dry, errMsg := requireWrite(m)
	if errMsg != "" {
		return mcp.ToolResultError(errMsg)
	}
	input := map[string]interface{}{"title": strArg(m, "title")}
	if d := strArg(m, "description"); d != "" {
		input["descriptionHtml"] = d
	}
	if v := strArg(m, "vendor"); v != "" {
		input["vendor"] = v
	}
	if pt := strArg(m, "product_type"); pt != "" {
		input["productType"] = pt
	}
	if st := strings.ToUpper(strArg(m, "status")); st != "" {
		input["status"] = st
	}
	if dry {
		return jsonResult(map[string]interface{}{"dry_run": true, "mutation": "productCreate", "input": input})
	}
	ctx, cancel := ctx()
	defer cancel()
	const q = `mutation($input: ProductInput!) {
	  productCreate(input: $input) {
	    product { id title status }
	    userErrors { field message }
	  }
	}`
	data, err := c.GraphQL(ctx, q, map[string]interface{}{"input": input})
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	if err := shopify.ExtractUserErrors(data, "productCreate", "userErrors"); err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return gqlOK(data, nil)
}

func handleUpdateProduct(raw json.RawMessage) map[string]interface{} {
	c, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	dry, errMsg := requireWrite(m)
	if errMsg != "" {
		return mcp.ToolResultError(errMsg)
	}
	id := productGID(strArg(m, "id"))
	input := map[string]interface{}{"id": id}
	if t := strArg(m, "title"); t != "" {
		input["title"] = t
	}
	if d := strArg(m, "description"); d != "" {
		input["descriptionHtml"] = d
	}
	if v := strArg(m, "vendor"); v != "" {
		input["vendor"] = v
	}
	if st := strings.ToUpper(strArg(m, "status")); st != "" {
		input["status"] = st
	}
	if dry {
		return jsonResult(map[string]interface{}{"dry_run": true, "mutation": "productUpdate", "input": input})
	}
	ctx, cancel := ctx()
	defer cancel()
	const q = `mutation($input: ProductInput!) {
	  productUpdate(input: $input) {
	    product { id title status }
	    userErrors { field message }
	  }
	}`
	data, err := c.GraphQL(ctx, q, map[string]interface{}{"input": input})
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	if err := shopify.ExtractUserErrors(data, "productUpdate", "userErrors"); err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return gqlOK(data, nil)
}

func handleListCollections(raw json.RawMessage) map[string]interface{} {
	c, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	ctx, cancel := ctx()
	defer cancel()
	vars := pageVars(m)
	const q = `query($first: Int!, $after: String, $query: String) {
	  collections(first: $first, after: $after, query: $query) {
	    pageInfo { hasNextPage endCursor }
	    edges { node { id title handle updatedAt } }
	  }
	}`
	data, err := c.GraphQL(ctx, q, vars)
	return gqlOK(data, err)
}

func handleGetCollection(raw json.RawMessage) map[string]interface{} {
	c, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	id := collectionGID(strArg(m, "id"))
	ctx, cancel := ctx()
	defer cancel()
	const q = `query($id: ID!) {
	  collection(id: $id) { id title handle descriptionHtml updatedAt productsCount { count } }
	}`
	data, err := c.GraphQL(ctx, q, map[string]interface{}{"id": id})
	return gqlOK(data, err)
}

func handleListOrders(raw json.RawMessage) map[string]interface{} {
	c, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	ctx, cancel := ctx()
	defer cancel()
	vars := pageVars(m)
	const q = `query($first: Int!, $after: String, $query: String) {
	  orders(first: $first, after: $after, query: $query, sortKey: CREATED_AT, reverse: true) {
	    pageInfo { hasNextPage endCursor }
	    edges { node { id name email displayFinancialStatus displayFulfillmentStatus createdAt totalPriceSet { shopMoney { amount currencyCode } } } }
	  }
	}`
	data, err := c.GraphQL(ctx, q, vars)
	return gqlOK(data, err)
}

func handleGetOrder(raw json.RawMessage) map[string]interface{} {
	c, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	id := orderGID(strArg(m, "id"))
	ctx, cancel := ctx()
	defer cancel()
	const q = `query($id: ID!) {
	  order(id: $id) {
	    id name email note tags createdAt displayFinancialStatus displayFulfillmentStatus
	    totalPriceSet { shopMoney { amount currencyCode } }
	    lineItems(first: 50) { edges { node { id title quantity sku variant { id } } } }
	    fulfillments { id status trackingInfo { number company url } }
	  }
	}`
	data, err := c.GraphQL(ctx, q, map[string]interface{}{"id": id})
	return gqlOK(data, err)
}

func handleUpdateOrder(raw json.RawMessage) map[string]interface{} {
	c, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	dry, errMsg := requireWrite(m)
	if errMsg != "" {
		return mcp.ToolResultError(errMsg)
	}
	id := orderGID(strArg(m, "id"))
	input := map[string]interface{}{"id": id}
	if note := strArg(m, "note"); note != "" {
		input["note"] = note
	}
	if tags, ok := m["tags"].([]interface{}); ok {
		tagStrs := make([]string, 0, len(tags))
		for _, t := range tags {
			tagStrs = append(tagStrs, fmt.Sprint(t))
		}
		input["tags"] = strings.Join(tagStrs, ", ")
	}
	if dry {
		return jsonResult(map[string]interface{}{"dry_run": true, "mutation": "orderUpdate", "input": input})
	}
	ctx, cancel := ctx()
	defer cancel()
	const q = `mutation($input: OrderInput!) {
	  orderUpdate(input: $input) {
	    order { id name note tags }
	    userErrors { field message }
	  }
	}`
	data, err := c.GraphQL(ctx, q, map[string]interface{}{"input": input})
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	if err := shopify.ExtractUserErrors(data, "orderUpdate", "userErrors"); err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return gqlOK(data, nil)
}

func handleCreateFulfillment(raw json.RawMessage) map[string]interface{} {
	c, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	dry, errMsg := requireWrite(m)
	if errMsg != "" {
		return mcp.ToolResultError(errMsg)
	}
	orderID := orderGID(strArg(m, "order_id"))
	notify := boolArg(m, "notify_customer")

	// Resolve fulfillment order + line items via GraphQL first unless dry_run of raw payload.
	planned := map[string]interface{}{
		"order_id":         orderID,
		"notify_customer":  notify,
		"tracking_number":  strArg(m, "tracking_number"),
		"tracking_company": strArg(m, "tracking_company"),
		"line_item_ids":    m["line_item_ids"],
	}
	if dry {
		return jsonResult(map[string]interface{}{"dry_run": true, "mutation": "fulfillmentCreateV2", "planned": planned})
	}

	ctx, cancel := ctx()
	defer cancel()
	const foQ = `query($id: ID!) {
	  order(id: $id) {
	    fulfillmentOrders(first: 10) {
	      edges { node {
	        id status
	        lineItems(first: 50) { edges { node { id remainingQuantity lineItem { id } } } }
	      } }
	    }
	  }
	}`
	foData, err := c.GraphQL(ctx, foQ, map[string]interface{}{"id": orderID})
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	var foParsed struct {
		Order struct {
			FulfillmentOrders struct {
				Edges []struct {
					Node struct {
						ID        string `json:"id"`
						Status    string `json:"status"`
						LineItems struct {
							Edges []struct {
								Node struct {
									ID                string `json:"id"`
									RemainingQuantity int    `json:"remainingQuantity"`
									LineItem          struct {
										ID string `json:"id"`
									} `json:"lineItem"`
								} `json:"node"`
							} `json:"edges"`
						} `json:"lineItems"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"fulfillmentOrders"`
		} `json:"order"`
	}
	if err := json.Unmarshal(foData, &foParsed); err != nil {
		return mcp.ToolResultError("parse fulfillment orders: " + err.Error())
	}

	wantLineItems := map[string]bool{}
	if arr, ok := m["line_item_ids"].([]interface{}); ok {
		for _, v := range arr {
			wantLineItems[fmt.Sprint(v)] = true
		}
	}

	var fulfillmentOrderID string
	lineItemsByFulfillmentOrder := []map[string]interface{}{}
	for _, e := range foParsed.Order.FulfillmentOrders.Edges {
		n := e.Node
		if n.Status != "OPEN" && n.Status != "IN_PROGRESS" {
			continue
		}
		items := []map[string]interface{}{}
		for _, li := range n.LineItems.Edges {
			if li.Node.RemainingQuantity <= 0 {
				continue
			}
			if len(wantLineItems) > 0 {
				if !wantLineItems[li.Node.LineItem.ID] && !wantLineItems[li.Node.ID] {
					continue
				}
			}
			items = append(items, map[string]interface{}{
				"id":       li.Node.ID,
				"quantity": li.Node.RemainingQuantity,
			})
		}
		if len(items) == 0 {
			continue
		}
		fulfillmentOrderID = n.ID
		lineItemsByFulfillmentOrder = append(lineItemsByFulfillmentOrder, map[string]interface{}{
			"fulfillmentOrderId":        n.ID,
			"fulfillmentOrderLineItems": items,
		})
		break
	}
	if fulfillmentOrderID == "" || len(lineItemsByFulfillmentOrder) == 0 {
		return mcp.ToolResultError("no open fulfillment order line items found for order")
	}

	tracking := map[string]interface{}{}
	if tn := strArg(m, "tracking_number"); tn != "" {
		tracking["number"] = tn
	}
	if tc := strArg(m, "tracking_company"); tc != "" {
		tracking["company"] = tc
	}

	fulfillment := map[string]interface{}{
		"lineItemsByFulfillmentOrder": lineItemsByFulfillmentOrder,
		"notifyCustomer":              notify,
	}
	if len(tracking) > 0 {
		fulfillment["trackingInfo"] = tracking
	}

	const mut = `mutation($fulfillment: FulfillmentV2Input!) {
	  fulfillmentCreateV2(fulfillment: $fulfillment) {
	    fulfillment { id status }
	    userErrors { field message }
	  }
	}`
	data, err := c.GraphQL(ctx, mut, map[string]interface{}{"fulfillment": fulfillment})
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	if err := shopify.ExtractUserErrors(data, "fulfillmentCreateV2", "userErrors"); err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return gqlOK(data, nil)
}

func handleListCustomers(raw json.RawMessage) map[string]interface{} {
	c, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	ctx, cancel := ctx()
	defer cancel()
	vars := pageVars(m)
	const q = `query($first: Int!, $after: String, $query: String) {
	  customers(first: $first, after: $after, query: $query) {
	    pageInfo { hasNextPage endCursor }
	    edges { node { id displayName email phone createdAt } }
	  }
	}`
	data, err := c.GraphQL(ctx, q, vars)
	return gqlOK(data, err)
}

func handleGetCustomer(raw json.RawMessage) map[string]interface{} {
	c, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	id := customerGID(strArg(m, "id"))
	ctx, cancel := ctx()
	defer cancel()
	const q = `query($id: ID!) {
	  customer(id: $id) {
	    id displayName email phone note tags createdAt
	    defaultAddress { address1 city province country zip }
	  }
	}`
	data, err := c.GraphQL(ctx, q, map[string]interface{}{"id": id})
	return gqlOK(data, err)
}

func handleCreateCustomer(raw json.RawMessage) map[string]interface{} {
	c, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	dry, errMsg := requireWrite(m)
	if errMsg != "" {
		return mcp.ToolResultError(errMsg)
	}
	input := map[string]interface{}{}
	if e := strArg(m, "email"); e != "" {
		input["email"] = e
	}
	if f := strArg(m, "first_name"); f != "" {
		input["firstName"] = f
	}
	if l := strArg(m, "last_name"); l != "" {
		input["lastName"] = l
	}
	if p := strArg(m, "phone"); p != "" {
		input["phone"] = p
	}
	if len(input) == 0 {
		return mcp.ToolResultError("at least one of email, first_name, last_name, phone required")
	}
	if dry {
		return jsonResult(map[string]interface{}{"dry_run": true, "mutation": "customerCreate", "input": input})
	}
	ctx, cancel := ctx()
	defer cancel()
	const q = `mutation($input: CustomerInput!) {
	  customerCreate(input: $input) {
	    customer { id email displayName }
	    userErrors { field message }
	  }
	}`
	data, err := c.GraphQL(ctx, q, map[string]interface{}{"input": input})
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	if err := shopify.ExtractUserErrors(data, "customerCreate", "userErrors"); err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return gqlOK(data, nil)
}

func handleUpdateCustomer(raw json.RawMessage) map[string]interface{} {
	c, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	dry, errMsg := requireWrite(m)
	if errMsg != "" {
		return mcp.ToolResultError(errMsg)
	}
	input := map[string]interface{}{"id": customerGID(strArg(m, "id"))}
	if e := strArg(m, "email"); e != "" {
		input["email"] = e
	}
	if f := strArg(m, "first_name"); f != "" {
		input["firstName"] = f
	}
	if l := strArg(m, "last_name"); l != "" {
		input["lastName"] = l
	}
	if p := strArg(m, "phone"); p != "" {
		input["phone"] = p
	}
	if dry {
		return jsonResult(map[string]interface{}{"dry_run": true, "mutation": "customerUpdate", "input": input})
	}
	ctx, cancel := ctx()
	defer cancel()
	const q = `mutation($input: CustomerInput!) {
	  customerUpdate(input: $input) {
	    customer { id email displayName }
	    userErrors { field message }
	  }
	}`
	data, err := c.GraphQL(ctx, q, map[string]interface{}{"input": input})
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	if err := shopify.ExtractUserErrors(data, "customerUpdate", "userErrors"); err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return gqlOK(data, nil)
}

func handleListInventoryLevels(raw json.RawMessage) map[string]interface{} {
	c, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	id := inventoryItemGID(strArg(m, "inventory_item_id"))
	first := intArg(m, "first", 25)
	ctx, cancel := ctx()
	defer cancel()
	const q = `query($id: ID!, $first: Int!) {
	  inventoryItem(id: $id) {
	    id sku tracked
	    inventoryLevels(first: $first) {
	      edges { node {
	        id
	        location { id name }
	        quantities(names: ["available","on_hand"]) { name quantity }
	      } }
	    }
	  }
	}`
	data, err := c.GraphQL(ctx, q, map[string]interface{}{"id": id, "first": first})
	return gqlOK(data, err)
}

func handleGetInventoryItem(raw json.RawMessage) map[string]interface{} {
	c, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	id := inventoryItemGID(strArg(m, "id"))
	ctx, cancel := ctx()
	defer cancel()
	const q = `query($id: ID!) {
	  inventoryItem(id: $id) { id sku tracked unitCost { amount currencyCode } }
	}`
	data, err := c.GraphQL(ctx, q, map[string]interface{}{"id": id})
	return gqlOK(data, err)
}

func handleAdjustInventory(raw json.RawMessage) map[string]interface{} {
	c, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	dry, errMsg := requireWrite(m)
	if errMsg != "" {
		return mcp.ToolResultError(errMsg)
	}
	itemID := inventoryItemGID(strArg(m, "inventory_item_id"))
	locID := locationGID(strArg(m, "location_id"))
	delta := intArg(m, "delta", 0)
	reason := strArg(m, "reason")
	if reason == "" {
		reason = "correction"
	}
	input := map[string]interface{}{
		"reason": reason,
		"name":   "available",
		"changes": []map[string]interface{}{
			{
				"inventoryItemId": itemID,
				"locationId":      locID,
				"delta":           delta,
			},
		},
	}
	if dry {
		return jsonResult(map[string]interface{}{"dry_run": true, "mutation": "inventoryAdjustQuantities", "input": input})
	}
	ctx, cancel := ctx()
	defer cancel()
	const q = `mutation($input: InventoryAdjustQuantitiesInput!) {
	  inventoryAdjustQuantities(input: $input) {
	    inventoryAdjustmentGroup { createdAt reason }
	    userErrors { field message }
	  }
	}`
	data, err := c.GraphQL(ctx, q, map[string]interface{}{"input": input})
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	if err := shopify.ExtractUserErrors(data, "inventoryAdjustQuantities", "userErrors"); err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return gqlOK(data, nil)
}

func handleListDiscounts(raw json.RawMessage) map[string]interface{} {
	c, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	ctx, cancel := ctx()
	defer cancel()
	vars := pageVars(m)
	const q = `query($first: Int!, $after: String, $query: String) {
	  codeDiscountNodes(first: $first, after: $after, query: $query) {
	    pageInfo { hasNextPage endCursor }
	    edges { node {
	      id
	      codeDiscount {
	        ... on DiscountCodeBasic { title status codes(first: 5) { edges { node { code } } } }
	        ... on DiscountCodeBxgy { title status }
	        ... on DiscountCodeFreeShipping { title status }
	      }
	    } }
	  }
	}`
	data, err := c.GraphQL(ctx, q, vars)
	return gqlOK(data, err)
}

func handleGetDiscount(raw json.RawMessage) map[string]interface{} {
	c, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	id := discountGID(strArg(m, "id"))
	ctx, cancel := ctx()
	defer cancel()
	const q = `query($id: ID!) {
	  codeDiscountNode(id: $id) {
	    id
	    codeDiscount {
	      ... on DiscountCodeBasic {
	        title status startsAt endsAt
	        codes(first: 20) { edges { node { code } } }
	        customerGets { value { ... on DiscountPercentage { percentage } ... on DiscountAmount { amount { amount currencyCode } } } }
	      }
	    }
	  }
	}`
	data, err := c.GraphQL(ctx, q, map[string]interface{}{"id": id})
	return gqlOK(data, err)
}

func handleCreateDiscount(raw json.RawMessage) map[string]interface{} {
	c, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	dry, errMsg := requireWrite(m)
	if errMsg != "" {
		return mcp.ToolResultError(errMsg)
	}
	title := strArg(m, "title")
	code := strArg(m, "code")
	startsAt := strArg(m, "starts_at")
	if startsAt == "" {
		startsAt = "2020-01-01T00:00:00Z"
	}

	customerGets := map[string]interface{}{
		"value":    map[string]interface{}{},
		"items":    map[string]interface{}{"all": true},
		"appliesOnOneTimePurchase": true,
	}
	if pct, ok := m["percentage"].(float64); ok && pct > 0 {
		customerGets["value"] = map[string]interface{}{"percentage": pct / 100.0}
	} else if amt := strArg(m, "amount"); amt != "" {
		customerGets["value"] = map[string]interface{}{
			"discountAmount": map[string]interface{}{
				"amount":            amt,
				"appliesOnEachItem": false,
			},
		}
	} else {
		return mcp.ToolResultError("amount or percentage required")
	}

	basic := map[string]interface{}{
		"title":       title,
		"code":        code,
		"startsAt":    startsAt,
		"customerGets": customerGets,
		"customerSelection": map[string]interface{}{"all": true},
	}
	if lim := intArg(m, "usage_limit", 0); lim > 0 {
		basic["usageLimit"] = lim
	}
	if dry {
		return jsonResult(map[string]interface{}{"dry_run": true, "mutation": "discountCodeBasicCreate", "basicCodeDiscount": basic})
	}
	ctx, cancel := ctx()
	defer cancel()
	const q = `mutation($basicCodeDiscount: DiscountCodeBasicInput!) {
	  discountCodeBasicCreate(basicCodeDiscount: $basicCodeDiscount) {
	    codeDiscountNode { id }
	    userErrors { field message }
	  }
	}`
	data, err := c.GraphQL(ctx, q, map[string]interface{}{"basicCodeDiscount": basic})
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	if err := shopify.ExtractUserErrors(data, "discountCodeBasicCreate", "userErrors"); err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return gqlOK(data, nil)
}

func handleUpdateDiscount(raw json.RawMessage) map[string]interface{} {
	c, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	dry, errMsg := requireWrite(m)
	if errMsg != "" {
		return mcp.ToolResultError(errMsg)
	}
	id := discountGID(strArg(m, "id"))
	basic := map[string]interface{}{}
	if t := strArg(m, "title"); t != "" {
		basic["title"] = t
	}
	if dry {
		return jsonResult(map[string]interface{}{"dry_run": true, "mutation": "discountCodeBasicUpdate", "id": id, "basicCodeDiscount": basic})
	}
	ctx, cancel := ctx()
	defer cancel()
	const q = `mutation($id: ID!, $basicCodeDiscount: DiscountCodeBasicInput!) {
	  discountCodeBasicUpdate(id: $id, basicCodeDiscount: $basicCodeDiscount) {
	    codeDiscountNode { id }
	    userErrors { field message }
	  }
	}`
	data, err := c.GraphQL(ctx, q, map[string]interface{}{"id": id, "basicCodeDiscount": basic})
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	if err := shopify.ExtractUserErrors(data, "discountCodeBasicUpdate", "userErrors"); err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return gqlOK(data, nil)
}

func handleDeleteDiscount(raw json.RawMessage) map[string]interface{} {
	c, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	dry, errMsg := requireWrite(m)
	if errMsg != "" {
		return mcp.ToolResultError(errMsg)
	}
	id := discountGID(strArg(m, "id"))
	if dry {
		return jsonResult(map[string]interface{}{"dry_run": true, "mutation": "discountCodeDelete", "id": id})
	}
	ctx, cancel := ctx()
	defer cancel()
	const q = `mutation($id: ID!) {
	  discountCodeDelete(id: $id) {
	    deletedCodeDiscountId
	    userErrors { field message }
	  }
	}`
	data, err := c.GraphQL(ctx, q, map[string]interface{}{"id": id})
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	if err := shopify.ExtractUserErrors(data, "discountCodeDelete", "userErrors"); err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return gqlOK(data, nil)
}
