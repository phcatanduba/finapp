package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"finapp/internal/model"
)

type webhookLogRepository struct {
	db *pgxpool.Pool
}

func NewWebhookLogRepository(db *pgxpool.Pool) WebhookLogRepository {
	return &webhookLogRepository{db: db}
}

func (r *webhookLogRepository) Create(ctx context.Context, log *model.WebhookLog) error {
	payload, err := json.Marshal(log.Payload)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx,
		`INSERT INTO webhook_logs (id, pluggy_item_id, event, payload, processed, error_message, received_at)
		 VALUES (uuid_generate_v4(), $1, $2, $3, $4, $5, NOW())`,
		log.PluggyItemID, log.Event, payload, log.Processed, log.ErrorMessage,
	)
	return err
}
