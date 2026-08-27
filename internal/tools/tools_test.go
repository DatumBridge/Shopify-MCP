package tools_test

import (
	"encoding/json"
	"testing"

	"github.com/datumbridge/shopify-mcp/internal/tools"
)

func TestRequireWriteViaCreateProduct(t *testing.T) {
	descs, handlers := tools.Register()
	if len(descs) < 20 {
		t.Fatalf("expected broad tool surface, got %d", len(descs))
	}
	h := handlers["shopify_create_product"]
	creds := `{"shop":"example.myshopify.com","token":"t"}`
	raw, _ := json.Marshal(map[string]interface{}{
		"credentials_json": creds,
		"title":            "X",
	})
	res := h(raw)
	if res["isError"] != true {
		t.Fatalf("expected confirm gate, got %v", res)
	}
	raw2, _ := json.Marshal(map[string]interface{}{
		"credentials_json": creds,
		"title":            "X",
		"confirm":          true,
		"dry_run":          true,
	})
	res2 := h(raw2)
	if res2["isError"] == true {
		t.Fatalf("expected dry_run success, got %v", res2)
	}
}
