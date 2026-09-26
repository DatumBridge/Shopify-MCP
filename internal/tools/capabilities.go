package tools

// capabilitiesForTool returns DatumBridge Tool Registry tags for tools/list `_meta`.
func capabilitiesForTool(name string) []string {
	switch name {
	case "shopify_list_products", "shopify_get_product", "shopify_create_product",
		"shopify_update_product", "shopify_list_collections", "shopify_get_collection":
		return []string{name, "catalog"}
	case "shopify_list_orders", "shopify_get_order", "shopify_update_order":
		return []string{name, "orders"}
	case "shopify_create_fulfillment":
		return []string{name, "orders", "fulfillment"}
	case "shopify_list_customers", "shopify_get_customer", "shopify_create_customer",
		"shopify_update_customer":
		return []string{name, "customers"}
	case "shopify_list_inventory_levels", "shopify_get_inventory_item", "shopify_adjust_inventory":
		return []string{name, "inventory"}
	case "shopify_list_discounts", "shopify_get_discount", "shopify_create_discount",
		"shopify_update_discount", "shopify_delete_discount":
		return []string{name, "discounts"}
	default:
		return []string{name, "shop"}
	}
}
