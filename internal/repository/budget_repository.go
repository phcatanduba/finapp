package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"finapp/internal/model"
)

type budgetRepository struct {
	db *pgxpool.Pool
}

func NewBudgetRepository(db *pgxpool.Pool) BudgetRepository {
	return &budgetRepository{db: db}
}

func (r *budgetRepository) Create(ctx context.Context, b *model.Budget) error {
	query := `
		INSERT INTO budgets (id, user_id, category_id, name, amount, period, start_date, end_date, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW(),NOW())
		RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query,
		b.ID, b.UserID, b.CategoryID, b.Name, b.Amount, b.Period, b.StartDate, b.EndDate,
	).Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
}

func (r *budgetRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.Budget, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, category_id, name, amount, period, start_date, end_date, created_at, updated_at
		 FROM budgets WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanAll(rows)
}

func (r *budgetRepository) FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Budget, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, user_id, category_id, name, amount, period, start_date, end_date, created_at, updated_at
		 FROM budgets WHERE id = $1 AND user_id = $2`, id, userID)
	return r.scanOne(row)
}

func (r *budgetRepository) Update(ctx context.Context, b *model.Budget) error {
	_, err := r.db.Exec(ctx,
		`UPDATE budgets SET category_id=$1, name=$2, amount=$3, period=$4, start_date=$5, end_date=$6, updated_at=NOW()
		 WHERE id=$7 AND user_id=$8`,
		b.CategoryID, b.Name, b.Amount, b.Period, b.StartDate, b.EndDate, b.ID, b.UserID,
	)
	return err
}

func (r *budgetRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM budgets WHERE id = $1 AND user_id = $2`, id, userID)
	return err
}

func (r *budgetRepository) scanAll(rows pgx.Rows) ([]model.Budget, error) {
	var budgets []model.Budget
	for rows.Next() {
		b, err := r.scanOne(rows)
		if err != nil {
			return nil, err
		}
		budgets = append(budgets, *b)
	}
	return budgets, rows.Err()
}

func (r *budgetRepository) scanOne(row interface{ Scan(...interface{}) error }) (*model.Budget, error) {
	var b model.Budget
	err := row.Scan(&b.ID, &b.UserID, &b.CategoryID, &b.Name, &b.Amount, &b.Period, &b.StartDate, &b.EndDate, &b.CreatedAt, &b.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan budget: %w", err)
	}
	return &b, nil
}
