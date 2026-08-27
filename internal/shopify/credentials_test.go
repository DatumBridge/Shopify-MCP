package shopify_test

import (
	"testing"

	"github.com/datumbridge/shopify-mcp/internal/shopify"
)

func TestNormalizeShop(t *testing.T) {
	cases := map[string]string{
		"Example":                        "example.myshopify.com",
		"https://Example.myshopify.com/": "example.myshopify.com",
		"example.myshopify.com":          "example.myshopify.com",
		"evil.example.com":               "",
		"https://evil.example.com":       "",
	}
	for in, want := range cases {
		if got := shopify.NormalizeShop(in); got != want {
			t.Fatalf("%q => %q want %q", in, got, want)
		}
	}
}

func TestParseCredentials(t *testing.T) {
	c, err := shopify.ParseCredentials(`{"shop":"acme","token":"tok","client_secret":"secret"}`, "")
	if err != nil {
		t.Fatal(err)
	}
	if c.Shop != "acme.myshopify.com" || c.Token != "tok" {
		t.Fatalf("%+v", c)
	}
	if c.ClientSecret != "" {
		t.Fatal("client_secret must be stripped")
	}
	_, err = shopify.ParseCredentials(`{"shop":"acme"}`, "")
	if err == nil {
		t.Fatal("expected missing token error")
	}
	_, err = shopify.ParseCredentials(`{"shop":"evil.example.com","token":"tok"}`, "")
	if err == nil {
		t.Fatal("expected reject non-myshopify host")
	}
}
