package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"finapp/internal/model"
	"finapp/internal/repository"
)

type projectionService struct {
	txRepo      repository.TransactionRepository
	accountRepo repository.AccountRepository
}

func NewProjectionService(txRepo repository.TransactionRepository, accountRepo repository.AccountRepository) ProjectionService {
	return &projectionService{txRepo: txRepo, accountRepo: accountRepo}
}

func (s *projectionService) ProjectBalance(ctx context.Context, userID uuid.UUID, months int) (*model.BalanceProjection, error) {
	if months <= 0 || months > 120 {
		return nil, fmt.Errorf("months must be between 1 and 120")
	}

	// Current total bank balance
	currentBalance, err := s.accountRepo.SumBalanceByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get current balance: %w", err)
	}

	// Compute monthly averages from last 3 months
	avgIncome, avgExpenses, err := s.txRepo.FindRecurring(ctx, userID, 3)
	if err != nil {
		return nil, fmt.Errorf("compute averages: %w", err)
	}

	points := make([]model.ProjectionPoint, months)
	balance := currentBalance
	now := time.Now()

	for i := 0; i < months; i++ {
		targetMonth := now.AddDate(0, i+1, 0)
		balance = balance + avgIncome - avgExpenses
		if balance < 0 {
			balance = 0
		}
		points[i] = model.ProjectionPoint{
			Month:            targetMonth.Format("2006-01"),
			ProjectedBalance: round2(balance),
			Income:           round2(avgIncome),
			Expenses:         round2(avgExpenses),
		}
	}

	return &model.BalanceProjection{
		CurrentBalance:    round2(currentBalance),
		MonthlyIncomeAvg:  round2(avgIncome),
		MonthlyExpenseAvg: round2(avgExpenses),
		ProjectedPoints:   points,
	}, nil
}
