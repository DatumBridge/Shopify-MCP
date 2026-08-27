package shopify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DefaultAPIVersion = "2025-01"

// Client talks to Shopify GraphQL Admin API.
type Client struct {
	HTTP       *http.Client
	Shop       string
	Token      string
	APIVersion string
}

func NewClient(creds *Credentials) *Client {
	ver := DefaultAPIVersion
	return &Client{
		HTTP: &http.Client{
			Timeout: 45 * time.Second,
		},
		Shop:       creds.Shop,
		Token:      creds.Token,
		APIVersion: ver,
	}
}

type gqlRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type gqlResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
	Extensions json.RawMessage `json:"extensions"`
}

// GraphQL executes a GraphQL query/mutation with bounded 429 retry.
func (c *Client) GraphQL(ctx context.Context, query string, variables map[string]interface{}) (json.RawMessage, error) {
	if c.Shop == "" || c.Token == "" {
		return nil, fmt.Errorf("shop and token required")
	}
	endpoint := fmt.Sprintf("https://%s/admin/api/%s/graphql.json", c.Shop, c.APIVersion)
	payload, err := json.Marshal(gqlRequest{Query: query, Variables: variables})
	if err != nil {
		return nil, err
	}

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("X-Shopify-Access-Token", c.Token)

		resp, err := c.HTTP.Do(req)
		if err != nil {
			return nil, err
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			wait := retryAfterSeconds(resp.Header.Get("Retry-After"), attempt)
			lastErr = fmt.Errorf("shopify rate limited (429)")
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(wait):
			}
			continue
		}
		if resp.StatusCode >= 300 {
			return nil, fmt.Errorf("shopify HTTP %d: %s", resp.StatusCode, truncate(string(body), 300))
		}

		var gr gqlResponse
		if err := json.Unmarshal(body, &gr); err != nil {
			return nil, fmt.Errorf("decode graphql response: %w", err)
		}
		if len(gr.Errors) > 0 {
			msgs := make([]string, 0, len(gr.Errors))
			for _, e := range gr.Errors {
				msgs = append(msgs, e.Message)
			}
			return nil, fmt.Errorf("graphql: %s", strings.Join(msgs, "; "))
		}
		return gr.Data, nil
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("shopify request failed")
}

func retryAfterSeconds(header string, attempt int) time.Duration {
	if header != "" {
		if sec, err := strconv.Atoi(strings.TrimSpace(header)); err == nil && sec > 0 {
			return time.Duration(sec) * time.Second
		}
	}
	return time.Duration(250*(attempt+1)) * time.Millisecond
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

// ExtractUserErrors reads common Shopify mutation userErrors arrays from data.
func ExtractUserErrors(data json.RawMessage, path ...string) error {
	var root interface{}
	if err := json.Unmarshal(data, &root); err != nil {
		return nil
	}
	cur := root
	for _, p := range path {
		m, ok := cur.(map[string]interface{})
		if !ok {
			return nil
		}
		cur = m[p]
	}
	arr, ok := cur.([]interface{})
	if !ok || len(arr) == 0 {
		return nil
	}
	msgs := make([]string, 0, len(arr))
	for _, item := range arr {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if msg, ok := m["message"].(string); ok && msg != "" {
			msgs = append(msgs, msg)
		}
	}
	if len(msgs) == 0 {
		return nil
	}
	return fmt.Errorf("userErrors: %s", strings.Join(msgs, "; "))
}
