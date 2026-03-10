package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"finapp/internal/model"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.Name,
	).Scan(&user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, email, password_hash, name, created_at, updated_at FROM users WHERE email = $1`
	row := r.db.QueryRow(ctx, query, email)
	return scanUser(row)
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `SELECT id, email, password_hash, name, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	return scanUser(row)
}

func scanUser(row pgx.Row) (*model.User, error) {
	var u model.User
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}
	return &u, nil
}
