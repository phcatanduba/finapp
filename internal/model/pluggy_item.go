package model

import (
	"time"

	"github.com/google/uuid"
)

type PluggyItem struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	PluggyItemID  string     `json:"pluggy_item_id"`
	ConnectorName string     `json:"connector_name"`
	ConnectorID   *int       `json:"connector_id,omitempty"`
	Status        string     `json:"status"`
	LastSyncedAt  *time.Time `json:"last_synced_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type ConnectTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type SyncRequest struct {
	ItemID string `json:"item_id,omitempty"` // if empty, sync all items
}
