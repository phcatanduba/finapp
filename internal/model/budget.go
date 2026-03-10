package model

import (
	"time"

	"github.com/google/uuid"
)

type BudgetPeriod string

const (
	BudgetPeriodMonthly BudgetPeriod = "MONTHLY"
	BudgetPeriodYearly  BudgetPeriod = "YEARLY"
)

type Budget struct {
	ID         uuid.UUID    `json:"id"`
	UserID     uuid.UUID    `json:"user_id"`
	CategoryID *uuid.UUID   `json:"category_id,omitempty"`
	Name       string       `json:"name"`
	Amount     float64      `json:"amount"`
	Period     BudgetPeriod `json:"period"`
	StartDate  time.Time    `json:"start_date"`
	EndDate    *time.Time   `json:"end_date,omitempty"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

type BudgetRequest struct {
	CategoryID *uuid.UUID   `json:"category_id"`
	Name       string       `json:"name"`
	Amount     float64      `json:"amount"`
	Period     BudgetPeriod `json:"period"`
	StartDate  time.Time    `json:"start_date"`
	EndDate    *time.Time   `json:"end_date"`
}

type BudgetProgress struct {
	Budget     Budget  `json:"budget"`
	Spent      float64 `json:"spent"`
	Remaining  float64 `json:"remaining"`
	Percentage float64 `json:"percentage"`
	PeriodFrom string  `json:"period_from"`
	PeriodTo   string  `json:"period_to"`
}
