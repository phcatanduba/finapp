package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"finapp/internal/model"
	"finapp/internal/repository"
)

type transactionService struct {
	txRepo repository.TransactionRepository
}

func NewTransactionService(txRepo repository.TransactionRepository) TransactionService {
	return &transactionService{txRepo: txRepo}
}

func (s *transactionService) ListTransactions(ctx context.Context, filter model.TransactionFilter) (*model.TransactionListResponse, error) {
	txs, total, err := s.txRepo.FindByFilter(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list transactions: %w", err)
	}
	if txs == nil {
		txs = []model.Transaction{}
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	return &model.TransactionListResponse{
		Transactions: txs,
		Total:        total,
		Page:         page,
		Limit:        limit,
	}, nil
}

func (s *transactionService) GetTransaction(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Transaction, error) {
	tx, err := s.txRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("get transaction: %w", err)
	}
	if tx == nil {
		return nil, ErrNotFound
	}
	return tx, nil
}

func (s *transactionService) UpdateTransaction(ctx context.Context, id uuid.UUID, userID uuid.UUID, req model.UpdateTransactionRequest) (*model.Transaction, error) {
	tx, err := s.txRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("find transaction: %w", err)
	}
	if tx == nil {
		return nil, ErrNotFound
	}

	tx.CategoryID = req.CategoryID
	tx.Notes = req.Notes
	if req.Tags != nil {
		tx.Tags = req.Tags
	}

	if err := s.txRepo.Update(ctx, tx); err != nil {
		return nil, fmt.Errorf("update transaction: %w", err)
	}
	return tx, nil
}
