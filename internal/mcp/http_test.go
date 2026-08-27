package mcp_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/datumbridge/shopify-mcp/internal/mcp"
	"github.com/datumbridge/shopify-mcp/internal/tools"
)

func newTestServer() *mcp.Server {
	descs, handlers := tools.Register()
	return &mcp.Server{
		Sessions: mcp.NewSessionStore(),
		Tools:    descs,
		Handlers: handlers,
	}
}

func postRPC(t *testing.T, s *mcp.Server, body map[string]interface{}, session string) *httptest.ResponseRecorder {
	t.Helper()
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/mcp", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	if session != "" {
		req.Header.Set(mcp.SessionHeader, session)
	}
	rr := httptest.NewRecorder()
	s.HandleStreamableHTTP(rr, req)
	return rr
}

func TestInitializeAndToolsList(t *testing.T) {
	s := newTestServer()
	rr := postRPC(t, s, map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params":  map[string]interface{}{},
	}, "")
	if rr.Code != 200 {
		t.Fatalf("status %d", rr.Code)
	}
	sid := rr.Header().Get(mcp.SessionHeader)
	if sid == "" {
		t.Fatal("expected Mcp-Session-Id")
	}

	rr2 := postRPC(t, s, map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/list",
		"params":  map[string]interface{}{},
	}, sid)
	var resp map[string]interface{}
	if err := json.Unmarshal(rr2.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	result, ok := resp["result"].(map[string]interface{})
	if !ok {
		t.Fatalf("no result: %s", rr2.Body.String())
	}
	toolsList, ok := result["tools"].([]interface{})
	if !ok || len(toolsList) < 10 {
		t.Fatalf("expected many tools, got %v", result["tools"])
	}
	capsByName := map[string]string{}
	for _, raw := range toolsList {
		tool, _ := raw.(map[string]interface{})
		name, _ := tool["name"].(string)
		meta, _ := tool["_meta"].(map[string]interface{})
		if meta == nil {
			t.Fatalf("tool %s missing _meta", name)
		}
		rawCaps, _ := meta["capabilities"].([]interface{})
		if len(rawCaps) == 0 {
			t.Fatalf("tool %s missing _meta.capabilities", name)
		}
		joined := ""
		for i, c := range rawCaps {
			if i > 0 {
				joined += ","
			}
			joined += c.(string)
		}
		capsByName[name] = joined
	}
	if capsByName["shopify_get_shop"] == capsByName["shopify_list_products"] {
		t.Fatalf("expected distinct capabilities, got %s", capsByName["shopify_get_shop"])
	}
	if capsByName["shopify_list_orders"] == capsByName["shopify_list_products"] {
		t.Fatalf("orders and catalog should differ")
	}
}

func TestToolsListRejectsMissingSession(t *testing.T) {
	s := newTestServer()
	rr := postRPC(t, s, map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/list",
	}, "")
	var resp map[string]interface{}
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	errObj, _ := resp["error"].(map[string]interface{})
	if errObj == nil {
		t.Fatal("expected error")
	}
	if int(errObj["code"].(float64)) != -32000 {
		t.Fatalf("code=%v", errObj["code"])
	}
}

func TestUnknownTool(t *testing.T) {
	s := newTestServer()
	rr := postRPC(t, s, map[string]interface{}{
		"jsonrpc": "2.0", "id": 1, "method": "initialize",
	}, "")
	sid := rr.Header().Get(mcp.SessionHeader)
	rr2 := postRPC(t, s, map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name":      "not_a_real_tool",
			"arguments": map[string]interface{}{},
		},
	}, sid)
	var resp map[string]interface{}
	_ = json.Unmarshal(rr2.Body.Bytes(), &resp)
	errObj, _ := resp["error"].(map[string]interface{})
	if errObj == nil || int(errObj["code"].(float64)) != -32601 {
		t.Fatalf("expected -32601, got %v", resp)
	}
}

func TestWriteRequiresConfirm(t *testing.T) {
	s := newTestServer()
	rr := postRPC(t, s, map[string]interface{}{
		"jsonrpc": "2.0", "id": 1, "method": "initialize",
	}, "")
	sid := rr.Header().Get(mcp.SessionHeader)
	creds := `{"shop":"example.myshopify.com","token":"shpat_test"}`
	rr2 := postRPC(t, s, map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "shopify_create_product",
			"arguments": map[string]interface{}{
				"credentials_json": creds,
				"title":            "Test",
			},
		},
	}, sid)
	var resp map[string]interface{}
	_ = json.Unmarshal(rr2.Body.Bytes(), &resp)
	result, _ := resp["result"].(map[string]interface{})
	if result == nil || result["isError"] != true {
		t.Fatalf("expected isError tool result, got %v", resp)
	}
}

func TestDryRunCreateProduct(t *testing.T) {
	s := newTestServer()
	rr := postRPC(t, s, map[string]interface{}{
		"jsonrpc": "2.0", "id": 1, "method": "initialize",
	}, "")
	sid := rr.Header().Get(mcp.SessionHeader)
	creds := `{"shop":"example.myshopify.com","token":"shpat_test"}`
	rr2 := postRPC(t, s, map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "shopify_create_product",
			"arguments": map[string]interface{}{
				"credentials_json": creds,
				"title":            "Test",
				"dry_run":          true,
			},
		},
	}, sid)
	var resp map[string]interface{}
	_ = json.Unmarshal(rr2.Body.Bytes(), &resp)
	result, _ := resp["result"].(map[string]interface{})
	if result == nil || result["isError"] == true {
		t.Fatalf("expected success dry_run, got %v", resp)
	}
	content, _ := result["content"].([]interface{})
	text := content[0].(map[string]interface{})["text"].(string)
	if !bytes.Contains([]byte(text), []byte(`"dry_run"`)) {
		t.Fatalf("expected dry_run payload, got %s", text)
	}
}

func TestToolsCallRejectsMissingSession(t *testing.T) {
	s := newTestServer()
	rr := postRPC(t, s, map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name":      "shopify_get_shop",
			"arguments": map[string]interface{}{},
		},
	}, "")
	var resp map[string]interface{}
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	errObj, _ := resp["error"].(map[string]interface{})
	if errObj == nil || int(errObj["code"].(float64)) != -32000 {
		t.Fatalf("expected -32000, got %v", resp)
	}
}

func TestCreateDiscountRequiresAmountOrPercentage(t *testing.T) {
	s := newTestServer()
	rr := postRPC(t, s, map[string]interface{}{
		"jsonrpc": "2.0", "id": 1, "method": "initialize",
	}, "")
	sid := rr.Header().Get(mcp.SessionHeader)
	creds := `{"shop":"example.myshopify.com","token":"shpat_test"}`
	rr2 := postRPC(t, s, map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "shopify_create_discount",
			"arguments": map[string]interface{}{
				"credentials_json": creds,
				"title":            "Sale",
				"code":             "SAVE",
				"dry_run":          true,
			},
		},
	}, sid)
	var resp map[string]interface{}
	_ = json.Unmarshal(rr2.Body.Bytes(), &resp)
	result, _ := resp["result"].(map[string]interface{})
	if result == nil || result["isError"] != true {
		t.Fatalf("expected isError for missing amount/percentage, got %v", resp)
	}
}
