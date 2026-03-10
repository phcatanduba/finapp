package model

import (
	"time"

	"github.com/google/uuid"
)

type WebhookPayload struct {
	Event     string                 `json:"event"`
	ItemID    string                 `json:"itemId"`
	Error     *WebhookError          `json:"error,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

type WebhookError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type WebhookLog struct {
	ID           uuid.UUID   `json:"id"`
	PluggyItemID *string     `json:"pluggy_item_id,omitempty"`
	Event        string      `json:"event"`
	Payload      interface{} `json:"payload"`
	Processed    bool        `json:"processed"`
	ErrorMessage *string     `json:"error_message,omitempty"`
	ReceivedAt   time.Time   `json:"received_at"`
}
