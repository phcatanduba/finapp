package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"finapp/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}

type PluggyItemRepository interface {
	Upsert(ctx context.Context, item *model.PluggyItem) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.PluggyItem, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.PluggyItem, error)
	FindByPluggyItemID(ctx context.Context, pluggyItemID string) (*model.PluggyItem, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	UpdateStatus(ctx context.Context, pluggyItemID string, status string, syncedAt *time.Time) error
}

type AccountRepository interface {
	Upsert(ctx context.Context, account *model.Account) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.Account, error)
	FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Account, error)
	FindByItemID(ctx context.Context, itemID uuid.UUID) ([]model.Account, error)
	SumBalanceByUserID(ctx context.Context, userID uuid.UUID) (float64, error)
}

type TransactionRepository interface {
	Upsert(ctx context.Context, tx *model.Transaction) error
	BulkUpsert(ctx context.Context, txs []model.Transaction) error
	FindByFilter(ctx context.Context, filter model.TransactionFilter) ([]model.Transaction, int, error)
	FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Transaction, error)
	Update(ctx context.Context, tx *model.Transaction) error
	SumByPeriod(ctx context.Context, userID uuid.UUID, from, to time.Time) (model.PeriodSummary, error)
	SumByCategoryAndPeriod(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]model.CategorySpending, error)
	GetMonthlyAggregates(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]model.CashFlowPoint, error)
	FindRecurring(ctx context.Context, userID uuid.UUID, lookbackMonths int) (float64, float64, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, cat *model.Category) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.Category, error)
	FindSystemCategories(ctx context.Context) ([]model.Category, error)
	HasSystemCategories(ctx context.Context) (bool, error)
	SeedSystemCategories(ctx context.Context, cats []model.SystemCategory) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Category, error)
	Update(ctx context.Context, cat *model.Category) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type BudgetRepository interface {
	Create(ctx context.Context, b *model.Budget) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.Budget, error)
	FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Budget, error)
	Update(ctx context.Context, b *model.Budget) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type GoalRepository interface {
	Create(ctx context.Context, g *model.Goal) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.Goal, error)
	FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Goal, error)
	Update(ctx context.Context, g *model.Goal) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type WebhookLogRepository interface {
	Create(ctx context.Context, log *model.WebhookLog) error
}
