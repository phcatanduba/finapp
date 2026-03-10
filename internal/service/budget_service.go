package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"finapp/internal/model"
	"finapp/internal/repository"
)

type budgetService struct {
	budgets  repository.BudgetRepository
	txRepo   repository.TransactionRepository
}

func NewBudgetService(budgets repository.BudgetRepository, txRepo repository.TransactionRepository) BudgetService {
	return &budgetService{budgets: budgets, txRepo: txRepo}
}

func (s *budgetService) Create(ctx context.Context, userID uuid.UUID, req model.BudgetRequest) (*model.Budget, error) {
	if req.Name == "" || req.Amount <= 0 {
		return nil, fmt.Errorf("name and amount are required")
	}
	b := &model.Budget{
		ID:         uuid.New(),
		UserID:     userID,
		CategoryID: req.CategoryID,
		Name:       req.Name,
		Amount:     req.Amount,
		Period:     req.Period,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
	}
	if err := s.budgets.Create(ctx, b); err != nil {
		return nil, fmt.Errorf("create budget: %w", err)
	}
	return b, nil
}

func (s *budgetService) List(ctx context.Context, userID uuid.UUID) ([]model.Budget, error) {
	return s.budgets.FindByUserID(ctx, userID)
}

func (s *budgetService) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Budget, error) {
	b, err := s.budgets.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, ErrNotFound
	}
	return b, nil
}

func (s *budgetService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req model.BudgetRequest) (*model.Budget, error) {
	b, err := s.budgets.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, ErrNotFound
	}

	b.CategoryID = req.CategoryID
	b.Name = req.Name
	b.Amount = req.Amount
	b.Period = req.Period
	b.StartDate = req.StartDate
	b.EndDate = req.EndDate

	if err := s.budgets.Update(ctx, b); err != nil {
		return nil, fmt.Errorf("update budget: %w", err)
	}
	return b, nil
}

func (s *budgetService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	b, err := s.budgets.FindByID(ctx, id, userID)
	if err != nil {
		return err
	}
	if b == nil {
		return ErrNotFound
	}
	return s.budgets.Delete(ctx, id, userID)
}

func (s *budgetService) GetProgress(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.BudgetProgress, error) {
	b, err := s.budgets.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, ErrNotFound
	}

	from, to := periodBounds(b.Period, time.Now())
	summary, err := s.txRepo.SumByPeriod(ctx, userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("sum transactions: %w", err)
	}

	spent := summary.TotalExpenses
	remaining := b.Amount - spent
	pct := 0.0
	if b.Amount > 0 {
		pct = (spent / b.Amount) * 100
	}

	return &model.BudgetProgress{
		Budget:     *b,
		Spent:      spent,
		Remaining:  remaining,
		Percentage: pct,
		PeriodFrom: from.Format("2006-01-02"),
		PeriodTo:   to.Format("2006-01-02"),
	}, nil
}

func periodBounds(period model.BudgetPeriod, now time.Time) (time.Time, time.Time) {
	switch period {
	case model.BudgetPeriodYearly:
		from := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		to := time.Date(now.Year(), 12, 31, 23, 59, 59, 0, now.Location())
		return from, to
	default: // MONTHLY
		from := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		to := from.AddDate(0, 1, -1)
		return from, to
	}
}
