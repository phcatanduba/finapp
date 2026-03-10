CREATE TABLE webhook_logs (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pluggy_item_id  VARCHAR(255),
    event           VARCHAR(100) NOT NULL,
    payload         JSONB NOT NULL,
    processed       BOOLEAN NOT NULL DEFAULT FALSE,
    error_message   TEXT,
    received_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_webhook_logs_pluggy_item_id ON webhook_logs(pluggy_item_id);
CREATE INDEX idx_webhook_logs_received_at ON webhook_logs(received_at DESC);
