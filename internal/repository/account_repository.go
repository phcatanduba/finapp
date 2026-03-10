package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"finapp/internal/model"
)

type accountRepository struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Upsert(ctx context.Context, acc *model.Account) error {
	query := `
		INSERT INTO accounts (id, user_id, item_id, pluggy_account_id, name, type, subtype, balance, credit_limit, currency_code, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		ON CONFLICT (pluggy_account_id) DO UPDATE
		SET name          = EXCLUDED.name,
		    balance       = EXCLUDED.balance,
		    credit_limit  = EXCLUDED.credit_limit,
		    updated_at    = NOW()
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		acc.ID, acc.UserID, acc.ItemID, acc.PluggyAccountID,
		acc.Name, acc.Type, acc.Subtype, acc.Balance, acc.CreditLimit, acc.CurrencyCode,
	).Scan(&acc.ID, &acc.CreatedAt, &acc.UpdatedAt)
}

func (r *accountRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.Account, error) {
	query := `
		SELECT id, user_id, item_id, pluggy_account_id, name, type, subtype, balance, credit_limit, currency_code, created_at, updated_at
		FROM accounts WHERE user_id = $1 ORDER BY name`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query accounts: %w", err)
	}
	defer rows.Close()
	return r.scanAll(rows)
}

func (r *accountRepository) FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Account, error) {
	query := `
		SELECT id, user_id, item_id, pluggy_account_id, name, type, subtype, balance, credit_limit, currency_code, created_at, updated_at
		FROM accounts WHERE id = $1 AND user_id = $2`
	return r.scanOne(r.db.QueryRow(ctx, query, id, userID))
}

func (r *accountRepository) FindByItemID(ctx context.Context, itemID uuid.UUID) ([]model.Account, error) {
	query := `
		SELECT id, user_id, item_id, pluggy_account_id, name, type, subtype, balance, credit_limit, currency_code, created_at, updated_at
		FROM accounts WHERE item_id = $1`
	rows, err := r.db.Query(ctx, query, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanAll(rows)
}

func (r *accountRepository) SumBalanceByUserID(ctx context.Context, userID uuid.UUID) (float64, error) {
	var total float64
	err := r.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(balance), 0) FROM accounts WHERE user_id = $1 AND type = 'BANK'`,
		userID,
	).Scan(&total)
	return total, err
}

func (r *accountRepository) scanAll(rows pgx.Rows) ([]model.Account, error) {
	var accounts []model.Account
	for rows.Next() {
		var a model.Account
		if err := r.scanRow(rows, &a); err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, rows.Err()
}

func (r *accountRepository) scanOne(row pgx.Row) (*model.Account, error) {
	var a model.Account
	err := row.Scan(&a.ID, &a.UserID, &a.ItemID, &a.PluggyAccountID,
		&a.Name, &a.Type, &a.Subtype, &a.Balance, &a.CreditLimit, &a.CurrencyCode,
		&a.CreatedAt, &a.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan account: %w", err)
	}
	return &a, nil
}

func (r *accountRepository) scanRow(rows pgx.Rows, a *model.Account) error {
	return rows.Scan(&a.ID, &a.UserID, &a.ItemID, &a.PluggyAccountID,
		&a.Name, &a.Type, &a.Subtype, &a.Balance, &a.CreditLimit, &a.CurrencyCode,
		&a.CreatedAt, &a.UpdatedAt)
}
