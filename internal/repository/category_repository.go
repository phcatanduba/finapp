package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"finapp/internal/model"
)

type categoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, cat *model.Category) error {
	query := `
		INSERT INTO categories (id, user_id, name, color, icon, is_system, parent_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		RETURNING created_at`
	return r.db.QueryRow(ctx, query,
		cat.ID, cat.UserID, cat.Name, cat.Color, cat.Icon, cat.IsSystem, cat.ParentID,
	).Scan(&cat.CreatedAt)
}

func (r *categoryRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.Category, error) {
	// returns system categories + user's own categories
	query := `
		SELECT id, user_id, name, color, icon, is_system, parent_id, created_at
		FROM categories
		WHERE is_system = TRUE OR user_id = $1
		ORDER BY is_system DESC, name`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanAll(rows)
}

func (r *categoryRepository) FindSystemCategories(ctx context.Context) ([]model.Category, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, name, color, icon, is_system, parent_id, created_at FROM categories WHERE is_system = TRUE`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanAll(rows)
}

func (r *categoryRepository) HasSystemCategories(ctx context.Context) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM categories WHERE is_system = TRUE`).Scan(&count)
	return count > 0, err
}

func (r *categoryRepository) SeedSystemCategories(ctx context.Context, cats []model.SystemCategory) error {
	for _, sc := range cats {
		_, err := r.db.Exec(ctx,
			`INSERT INTO categories (id, user_id, name, color, icon, is_system, created_at)
			 VALUES (uuid_generate_v4(), NULL, $1, $2, $3, TRUE, NOW())
			 ON CONFLICT DO NOTHING`,
			sc.Name, sc.Color, sc.Icon,
		)
		if err != nil {
			return fmt.Errorf("seed category %q: %w", sc.Name, err)
		}
	}
	return nil
}

func (r *categoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	query := `SELECT id, user_id, name, color, icon, is_system, parent_id, created_at FROM categories WHERE id = $1`
	return r.scanOne(r.db.QueryRow(ctx, query, id))
}

func (r *categoryRepository) Update(ctx context.Context, cat *model.Category) error {
	_, err := r.db.Exec(ctx,
		`UPDATE categories SET name=$1, color=$2, icon=$3, parent_id=$4 WHERE id=$5 AND user_id=$6`,
		cat.Name, cat.Color, cat.Icon, cat.ParentID, cat.ID, cat.UserID,
	)
	return err
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM categories WHERE id = $1 AND user_id = $2 AND is_system = FALSE`,
		id, userID,
	)
	return err
}

func (r *categoryRepository) scanAll(rows pgx.Rows) ([]model.Category, error) {
	var cats []model.Category
	for rows.Next() {
		cat, err := r.scanOne(rows)
		if err != nil {
			return nil, err
		}
		cats = append(cats, *cat)
	}
	return cats, rows.Err()
}

func (r *categoryRepository) scanOne(row interface{ Scan(...interface{}) error }) (*model.Category, error) {
	var cat model.Category
	err := row.Scan(&cat.ID, &cat.UserID, &cat.Name, &cat.Color, &cat.Icon, &cat.IsSystem, &cat.ParentID, &cat.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan category: %w", err)
	}
	return &cat, nil
}
