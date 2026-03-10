package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"finapp/internal/model"
)

type AuthService interface {
	Register(ctx context.Context, req model.RegisterRequest) (*model.AuthResponse, error)
	Login(ctx context.Context, req model.LoginRequest) (*model.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*model.AuthResponse, error)
	ValidateToken(tokenString string) (*model.Claims, error)
}

type PluggySyncService interface {
	GenerateConnectToken(ctx context.Context, userID uuid.UUID) (*model.ConnectTokenResponse, error)
	ListItems(ctx context.Context, userID uuid.UUID) ([]model.PluggyItem, error)
	DisconnectItem(ctx context.Context, itemID uuid.UUID, userID uuid.UUID) error
	SyncItem(ctx context.Context, pluggyItemID string, userID uuid.UUID) error
	SyncAllItems(ctx context.Context, userID uuid.UUID) error
	HandleWebhook(ctx context.Context, payload model.WebhookPayload) error
}

type AccountService interface {
	ListAccounts(ctx context.Context, userID uuid.UUID) ([]model.Account, error)
	GetAccount(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Account, error)
}

type TransactionService interface {
	ListTransactions(ctx context.Context, filter model.TransactionFilter) (*model.TransactionListResponse, error)
	GetTransaction(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Transaction, error)
	UpdateTransaction(ctx context.Context, id uuid.UUID, userID uuid.UUID, req model.UpdateTransactionRequest) (*model.Transaction, error)
}

type CategoryService interface {
	Create(ctx context.Context, userID uuid.UUID, req model.CategoryRequest) (*model.Category, error)
	List(ctx context.Context, userID uuid.UUID) ([]model.Category, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Category, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req model.CategoryRequest) (*model.Category, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type BudgetService interface {
	Create(ctx context.Context, userID uuid.UUID, req model.BudgetRequest) (*model.Budget, error)
	List(ctx context.Context, userID uuid.UUID) ([]model.Budget, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Budget, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req model.BudgetRequest) (*model.Budget, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetProgress(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.BudgetProgress, error)
}

type GoalService interface {
	Create(ctx context.Context, userID uuid.UUID, req model.GoalRequest) (*model.Goal, error)
	List(ctx context.Context, userID uuid.UUID) ([]model.Goal, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Goal, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req model.GoalRequest) (*model.Goal, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetProgress(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.GoalProgress, error)
}

type ReportService interface {
	GetSummary(ctx context.Context, userID uuid.UUID, from, to time.Time) (*model.ReportSummary, error)
	GetByCategory(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]model.CategorySpending, error)
	GetCashFlow(ctx context.Context, userID uuid.UUID, from, to time.Time) (*model.CashFlow, error)
}

type SimulationService interface {
	CompoundInterest(req model.CompoundInterestRequest) model.CompoundInterestResponse
	Loan(req model.LoanRequest) model.LoanResponse
}

type ProjectionService interface {
	ProjectBalance(ctx context.Context, userID uuid.UUID, months int) (*model.BalanceProjection, error)
}
