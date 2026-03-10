package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"finapp/internal/model"
)

type goalRepository struct {
	db *pgxpool.Pool
}

func NewGoalRepository(db *pgxpool.Pool) GoalRepository {
	return &goalRepository{db: db}
}

func (r *goalRepository) Create(ctx context.Context, g *model.Goal) error {
	query := `
		INSERT INTO goals (id, user_id, name, description, target_amount, current_amount, deadline, type, is_completed, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,FALSE,NOW(),NOW())
		RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query,
		g.ID, g.UserID, g.Name, g.Description, g.TargetAmount, g.CurrentAmount, g.Deadline, g.Type,
	).Scan(&g.ID, &g.CreatedAt, &g.UpdatedAt)
}

func (r *goalRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.Goal, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, name, description, target_amount, current_amount, deadline, type, is_completed, created_at, updated_at
		 FROM goals WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanAll(rows)
}

func (r *goalRepository) FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Goal, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, user_id, name, description, target_amount, current_amount, deadline, type, is_completed, created_at, updated_at
		 FROM goals WHERE id = $1 AND user_id = $2`, id, userID)
	return r.scanOne(row)
}

func (r *goalRepository) Update(ctx context.Context, g *model.Goal) error {
	_, err := r.db.Exec(ctx,
		`UPDATE goals SET name=$1, description=$2, target_amount=$3, current_amount=$4, deadline=$5, type=$6, is_completed=$7, updated_at=NOW()
		 WHERE id=$8 AND user_id=$9`,
		g.Name, g.Description, g.TargetAmount, g.CurrentAmount, g.Deadline, g.Type, g.IsCompleted, g.ID, g.UserID,
	)
	return err
}

func (r *goalRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM goals WHERE id = $1 AND user_id = $2`, id, userID)
	return err
}

func (r *goalRepository) scanAll(rows pgx.Rows) ([]model.Goal, error) {
	var goals []model.Goal
	for rows.Next() {
		g, err := r.scanOne(rows)
		if err != nil {
			return nil, err
		}
		goals = append(goals, *g)
	}
	return goals, rows.Err()
}

func (r *goalRepository) scanOne(row interface{ Scan(...interface{}) error }) (*model.Goal, error) {
	var g model.Goal
	err := row.Scan(&g.ID, &g.UserID, &g.Name, &g.Description, &g.TargetAmount, &g.CurrentAmount,
		&g.Deadline, &g.Type, &g.IsCompleted, &g.CreatedAt, &g.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan goal: %w", err)
	}
	return &g, nil
}
