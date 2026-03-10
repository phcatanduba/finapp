package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"finapp/internal/model"
	"finapp/internal/repository"
)

type reportService struct {
	txRepo repository.TransactionRepository
}

func NewReportService(txRepo repository.TransactionRepository) ReportService {
	return &reportService{txRepo: txRepo}
}

func (s *reportService) GetSummary(ctx context.Context, userID uuid.UUID, from, to time.Time) (*model.ReportSummary, error) {
	summary, err := s.txRepo.SumByPeriod(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}
	return &model.ReportSummary{
		TotalIncome:   summary.TotalIncome,
		TotalExpenses: summary.TotalExpenses,
		NetBalance:    summary.TotalIncome - summary.TotalExpenses,
		PeriodFrom:    from.Format("2006-01-02"),
		PeriodTo:      to.Format("2006-01-02"),
	}, nil
}

func (s *reportService) GetByCategory(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]model.CategorySpending, error) {
	return s.txRepo.SumByCategoryAndPeriod(ctx, userID, from, to)
}

func (s *reportService) GetCashFlow(ctx context.Context, userID uuid.UUID, from, to time.Time) (*model.CashFlow, error) {
	points, err := s.txRepo.GetMonthlyAggregates(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}
	if points == nil {
		points = []model.CashFlowPoint{}
	}
	return &model.CashFlow{
		Points:     points,
		PeriodFrom: from.Format("2006-01-02"),
		PeriodTo:   to.Format("2006-01-02"),
	}, nil
}
