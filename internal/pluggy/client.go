package pluggy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is the Pluggy API HTTP client.
type Client struct {
	baseURL      string
	clientID     string
	clientSecret string
	httpClient   *http.Client
	tokens       *tokenManager
}

func NewClient(baseURL, clientID, clientSecret string) *Client {
	c := &Client{
		baseURL:      baseURL,
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	c.tokens = newTokenManager(c)
	return c
}

// doRequest executes an authenticated HTTP request against the Pluggy API.
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	apiKey, err := c.tokens.GetAPIKey(ctx)
	if err != nil {
		return fmt.Errorf("get api key: %w", err)
	}
	return c.doRequestWithKey(ctx, method, path, apiKey, body, result)
}

// doRequestWithKey executes an HTTP request with an explicit API key.
func (c *Client) doRequestWithKey(ctx context.Context, method, path, apiKey string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pluggy API error %d: %s", resp.StatusCode, string(respBytes))
	}

	if result != nil {
		if err := json.Unmarshal(respBytes, result); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

// authenticate calls POST /auth and returns the raw API key.
func (c *Client) authenticate(ctx context.Context) (string, error) {
	payload := map[string]string{
		"clientId":     c.clientID,
		"clientSecret": c.clientSecret,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/auth", &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("authenticate: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("authenticate error %d: %s", resp.StatusCode, string(body))
	}

	var ar AuthResponse
	if err := json.Unmarshal(body, &ar); err != nil {
		return "", fmt.Errorf("decode auth response: %w", err)
	}

	return ar.APIKey, nil
}

// GenerateConnectToken returns a short-lived token for the Pluggy frontend widget.
func (c *Client) GenerateConnectToken(ctx context.Context, options map[string]interface{}) (*ConnectTokenResponse, error) {
	var result ConnectTokenResponse
	if err := c.doRequest(ctx, http.MethodPost, "/connect_token", options, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetItem retrieves an Item by ID.
func (c *Client) GetItem(ctx context.Context, itemID string) (*Item, error) {
	var result Item
	if err := c.doRequest(ctx, http.MethodGet, "/items/"+itemID, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAccounts returns all accounts linked to an Item.
func (c *Client) GetAccounts(ctx context.Context, itemID string) (*AccountsResponse, error) {
	var result AccountsResponse
	if err := c.doRequest(ctx, http.MethodGet, "/accounts?itemId="+itemID, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTransactions returns transactions for a given account with basic pagination.
func (c *Client) GetTransactions(ctx context.Context, accountID string, page int) (*TransactionsResponse, error) {
	path := fmt.Sprintf("/transactions?accountId=%s&pageSize=500&page=%d", accountID, page)
	var result TransactionsResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteItem disconnects an Item from Pluggy.
func (c *Client) DeleteItem(ctx context.Context, itemID string) error {
	return c.doRequest(ctx, http.MethodDelete, "/items/"+itemID, nil, nil)
}
