package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"finapp/internal/model"
	"finapp/internal/repository"
)

type accountService struct {
	accounts repository.AccountRepository
}

func NewAccountService(accounts repository.AccountRepository) AccountService {
	return &accountService{accounts: accounts}
}

func (s *accountService) ListAccounts(ctx context.Context, userID uuid.UUID) ([]model.Account, error) {
	return s.accounts.FindByUserID(ctx, userID)
}

func (s *accountService) GetAccount(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Account, error) {
	acc, err := s.accounts.FindByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("get account: %w", err)
	}
	if acc == nil {
		return nil, ErrNotFound
	}
	return acc, nil
}
