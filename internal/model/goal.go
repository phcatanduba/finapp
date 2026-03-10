package model

import (
	"time"

	"github.com/google/uuid"
)

type GoalType string

const (
	GoalTypeSavings       GoalType = "SAVINGS"
	GoalTypeDebtPayoff    GoalType = "DEBT_PAYOFF"
	GoalTypeInvestment    GoalType = "INVESTMENT"
	GoalTypeEmergencyFund GoalType = "EMERGENCY_FUND"
	GoalTypeOther         GoalType = "OTHER"
)

type Goal struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	Name          string     `json:"name"`
	Description   *string    `json:"description,omitempty"`
	TargetAmount  float64    `json:"target_amount"`
	CurrentAmount float64    `json:"current_amount"`
	Deadline      *time.Time `json:"deadline,omitempty"`
	Type          GoalType   `json:"type"`
	IsCompleted   bool       `json:"is_completed"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type GoalRequest struct {
	Name          string     `json:"name"`
	Description   *string    `json:"description"`
	TargetAmount  float64    `json:"target_amount"`
	CurrentAmount float64    `json:"current_amount"`
	Deadline      *time.Time `json:"deadline"`
	Type          GoalType   `json:"type"`
}

type GoalProgress struct {
	Goal          Goal    `json:"goal"`
	Remaining     float64 `json:"remaining"`
	Percentage    float64 `json:"percentage"`
	DaysToDeadline *int   `json:"days_to_deadline,omitempty"`
	OnTrack       bool    `json:"on_track"`
}
