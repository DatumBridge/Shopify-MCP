package shopify

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Credentials is the vault/injected credentials_json shape for Shopify Admin API.
type Credentials struct {
	Type         string   `json:"type"`
	Shop         string   `json:"shop"`
	Token        string   `json:"token"`
	AccessToken  string   `json:"access_token"` // alias
	RefreshToken string   `json:"refresh_token"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Scopes       []string `json:"scopes"`
	TokenURI     string   `json:"token_uri"`
}

// ParseCredentials loads Credentials from credentials_json string and/or credentials_path file.
func ParseCredentials(credentialsJSON, credentialsPath string) (*Credentials, error) {
	raw := strings.TrimSpace(credentialsJSON)
	if raw == "" && strings.TrimSpace(credentialsPath) != "" {
		path := filepath.Clean(credentialsPath)
		if strings.Contains(path, "..") {
			return nil, fmt.Errorf("credentials_path must not contain ..")
		}
		allowedRoot := strings.TrimSpace(os.Getenv("SHOPIFY_CREDENTIALS_DIR"))
		if allowedRoot == "" {
			allowedRoot = "/credentials"
		}
		allowedRoot = filepath.Clean(allowedRoot)
		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("credentials_path: %w", err)
		}
		absRoot, err := filepath.Abs(allowedRoot)
		if err != nil {
			return nil, fmt.Errorf("SHOPIFY_CREDENTIALS_DIR: %w", err)
		}
		if absPath != absRoot && !strings.HasPrefix(absPath, absRoot+string(os.PathSeparator)) {
			return nil, fmt.Errorf("credentials_path must be under %s", absRoot)
		}
		b, err := os.ReadFile(absPath)
		if err != nil {
			return nil, fmt.Errorf("read credentials_path: %w", err)
		}
		raw = strings.TrimSpace(string(b))
	}
	if raw == "" {
		return nil, fmt.Errorf("credentials_json or credentials_path required")
	}
	var c Credentials
	if err := json.Unmarshal([]byte(raw), &c); err != nil {
		return nil, fmt.Errorf("invalid credentials_json: %w", err)
	}
	c.Shop = NormalizeShop(c.Shop)
	token := strings.TrimSpace(c.Token)
	if token == "" {
		token = strings.TrimSpace(c.AccessToken)
	}
	c.Token = token
	if c.Shop == "" {
		return nil, fmt.Errorf("credentials missing shop (must be *.myshopify.com)")
	}
	if c.Token == "" {
		return nil, fmt.Errorf("credentials missing token / access_token")
	}
	// Never retain Partner app secret in-process from tool args.
	c.ClientSecret = ""
	return &c, nil
}

// NormalizeShop ensures a myshopify.com host (no scheme/path). Rejects non-Shopify hosts.
func NormalizeShop(shop string) string {
	s := strings.TrimSpace(strings.ToLower(shop))
	s = strings.TrimPrefix(s, "https://")
	s = strings.TrimPrefix(s, "http://")
	s = strings.TrimSuffix(s, "/")
	if i := strings.Index(s, "/"); i >= 0 {
		s = s[:i]
	}
	if s == "" {
		return ""
	}
	if !strings.Contains(s, ".") {
		s = s + ".myshopify.com"
	}
	if !strings.HasSuffix(s, ".myshopify.com") {
		return ""
	}
	return s
}
