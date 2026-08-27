package tools

import "testing"

func TestCapabilitiesForTool_distinctFamilies(t *testing.T) {
	shop := capabilitiesForTool("shopify_get_shop")
	product := capabilitiesForTool("shopify_list_products")
	order := capabilitiesForTool("shopify_list_orders")
	if containsCap(shop, "catalog") {
		t.Fatalf("shop identity should not be catalog: %v", shop)
	}
	if !containsCap(product, "catalog") {
		t.Fatalf("product tool missing catalog: %v", product)
	}
	if !containsCap(order, "orders") {
		t.Fatalf("order tool missing orders: %v", order)
	}
}

func containsCap(in []string, want string) bool {
	for _, s := range in {
		if s == want {
			return true
		}
	}
	return false
}
