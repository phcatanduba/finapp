package model

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeDebit  TransactionType = "DEBIT"
	TransactionTypeCredit TransactionType = "CREDIT"
)

type Transaction struct {
	ID                   uuid.UUID       `json:"id"`
	UserID               uuid.UUID       `json:"user_id"`
	AccountID            uuid.UUID       `json:"account_id"`
	PluggyTransactionID  *string         `json:"pluggy_transaction_id,omitempty"`
	Description          string          `json:"description"`
	Amount               float64         `json:"amount"`
	Date                 time.Time       `json:"date"`
	Type                 TransactionType `json:"type"`
	CategoryID           *uuid.UUID      `json:"category_id,omitempty"`
	PluggyCategory       *string         `json:"pluggy_category,omitempty"`
	Notes                *string         `json:"notes,omitempty"`
	Tags                 []string        `json:"tags"`
	IsRecurring          bool            `json:"is_recurring"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

type TransactionFilter struct {
	UserID     uuid.UUID
	AccountID  *uuid.UUID
	CategoryID *uuid.UUID
	Type       *TransactionType
	From       *time.Time
	To         *time.Time
	Page       int
	Limit      int
}

type UpdateTransactionRequest struct {
	CategoryID *uuid.UUID `json:"category_id"`
	Notes      *string    `json:"notes"`
	Tags       []string   `json:"tags"`
}

type TransactionListResponse struct {
	Transactions []Transaction `json:"transactions"`
	Total        int           `json:"total"`
	Page         int           `json:"page"`
	Limit        int           `json:"limit"`
}
