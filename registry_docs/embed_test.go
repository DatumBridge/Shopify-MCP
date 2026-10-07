package registrydocs

import (
	"strings"
	"testing"
)

func TestMarkdown_createProductIncludesConfirm(t *testing.T) {
	text := Markdown("shopify_create_product")
	if !strings.Contains(text, "dry_run") || !strings.Contains(text, "confirm") {
		t.Fatalf("missing write cases")
	}
	if Markdown("../shopify_create_product") != "" {
		t.Fatal("path escape returned a guide")
	}
}
