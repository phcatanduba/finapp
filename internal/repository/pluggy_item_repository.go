package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"finapp/internal/model"
)

type pluggyItemRepository struct {
	db *pgxpool.Pool
}

func NewPluggyItemRepository(db *pgxpool.Pool) PluggyItemRepository {
	return &pluggyItemRepository{db: db}
}

func (r *pluggyItemRepository) Upsert(ctx context.Context, item *model.PluggyItem) error {
	query := `
		INSERT INTO pluggy_items (id, user_id, pluggy_item_id, connector_name, connector_id, status, last_synced_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		ON CONFLICT (pluggy_item_id) DO UPDATE
		SET connector_name = EXCLUDED.connector_name,
		    connector_id   = EXCLUDED.connector_id,
		    status         = EXCLUDED.status,
		    last_synced_at = EXCLUDED.last_synced_at,
		    updated_at     = NOW()
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		item.ID, item.UserID, item.PluggyItemID, item.ConnectorName,
		item.ConnectorID, item.Status, item.LastSyncedAt,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
}

func (r *pluggyItemRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.PluggyItem, error) {
	query := `
		SELECT id, user_id, pluggy_item_id, connector_name, connector_id, status, last_synced_at, created_at, updated_at
		FROM pluggy_items WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query pluggy items: %w", err)
	}
	defer rows.Close()

	var items []model.PluggyItem
	for rows.Next() {
		var it model.PluggyItem
		if err := rows.Scan(&it.ID, &it.UserID, &it.PluggyItemID, &it.ConnectorName,
			&it.ConnectorID, &it.Status, &it.LastSyncedAt, &it.CreatedAt, &it.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan pluggy item: %w", err)
		}
		items = append(items, it)
	}
	return items, rows.Err()
}

func (r *pluggyItemRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.PluggyItem, error) {
	query := `
		SELECT id, user_id, pluggy_item_id, connector_name, connector_id, status, last_synced_at, created_at, updated_at
		FROM pluggy_items WHERE id = $1`
	return r.scanOne(r.db.QueryRow(ctx, query, id))
}

func (r *pluggyItemRepository) FindByPluggyItemID(ctx context.Context, pluggyItemID string) (*model.PluggyItem, error) {
	query := `
		SELECT id, user_id, pluggy_item_id, connector_name, connector_id, status, last_synced_at, created_at, updated_at
		FROM pluggy_items WHERE pluggy_item_id = $1`
	return r.scanOne(r.db.QueryRow(ctx, query, pluggyItemID))
}

func (r *pluggyItemRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM pluggy_items WHERE id = $1 AND user_id = $2`, id, userID)
	return err
}

func (r *pluggyItemRepository) UpdateStatus(ctx context.Context, pluggyItemID string, status string, syncedAt *time.Time) error {
	_, err := r.db.Exec(ctx,
		`UPDATE pluggy_items SET status = $1, last_synced_at = $2, updated_at = NOW() WHERE pluggy_item_id = $3`,
		status, syncedAt, pluggyItemID,
	)
	return err
}

func (r *pluggyItemRepository) scanOne(row pgx.Row) (*model.PluggyItem, error) {
	var it model.PluggyItem
	err := row.Scan(&it.ID, &it.UserID, &it.PluggyItemID, &it.ConnectorName,
		&it.ConnectorID, &it.Status, &it.LastSyncedAt, &it.CreatedAt, &it.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan pluggy item: %w", err)
	}
	return &it, nil
}
