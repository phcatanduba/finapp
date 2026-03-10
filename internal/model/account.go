package model

import (
	"time"

	"github.com/google/uuid"
)

type AccountType string

const (
	AccountTypeBank       AccountType = "BANK"
	AccountTypeCredit     AccountType = "CREDIT"
	AccountTypeInvestment AccountType = "INVESTMENT"
	AccountTypeLoan       AccountType = "LOAN"
	AccountTypeOther      AccountType = "OTHER"
)

type Account struct {
	ID               uuid.UUID   `json:"id"`
	UserID           uuid.UUID   `json:"user_id"`
	ItemID           uuid.UUID   `json:"item_id"`
	PluggyAccountID  string      `json:"pluggy_account_id"`
	Name             string      `json:"name"`
	Type             AccountType `json:"type"`
	Subtype          *string     `json:"subtype,omitempty"`
	Balance          float64     `json:"balance"`
	CreditLimit      *float64    `json:"credit_limit,omitempty"`
	CurrencyCode     string      `json:"currency_code"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
}
