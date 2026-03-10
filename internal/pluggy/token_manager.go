package pluggy

import (
	"context"
	"sync"
	"time"
)

// tokenManager caches the Pluggy API key which expires after 2 hours.
type tokenManager struct {
	mu        sync.Mutex
	apiKey    string
	expiresAt time.Time
	client    *Client
}

func newTokenManager(client *Client) *tokenManager {
	return &tokenManager{client: client}
}

// GetAPIKey returns a valid API key, refreshing it if expired or about to expire.
func (tm *tokenManager) GetAPIKey(ctx context.Context) (string, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Refresh if we have less than 5 minutes remaining
	if tm.apiKey == "" || time.Now().Add(5*time.Minute).After(tm.expiresAt) {
		if err := tm.refresh(ctx); err != nil {
			return "", err
		}
	}

	return tm.apiKey, nil
}

func (tm *tokenManager) refresh(ctx context.Context) error {
	apiKey, err := tm.client.authenticate(ctx)
	if err != nil {
		return err
	}
	tm.apiKey = apiKey
	tm.expiresAt = time.Now().Add(2 * time.Hour)
	return nil
}
