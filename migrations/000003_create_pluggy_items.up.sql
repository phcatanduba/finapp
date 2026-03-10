CREATE TABLE pluggy_items (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pluggy_item_id  VARCHAR(255) NOT NULL UNIQUE,
    connector_name  VARCHAR(255) NOT NULL,
    connector_id    INTEGER,
    status          VARCHAR(50) NOT NULL DEFAULT 'UPDATING',
    last_synced_at  TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pluggy_items_user_id ON pluggy_items(user_id);
CREATE INDEX idx_pluggy_items_pluggy_item_id ON pluggy_items(pluggy_item_id);
