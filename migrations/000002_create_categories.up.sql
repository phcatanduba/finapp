CREATE TABLE categories (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID REFERENCES users(id) ON DELETE CASCADE,
    name        VARCHAR(100) NOT NULL,
    color       VARCHAR(7),
    icon        VARCHAR(50),
    is_system   BOOLEAN NOT NULL DEFAULT FALSE,
    parent_id   UUID REFERENCES categories(id) ON DELETE SET NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_categories_user_id ON categories(user_id);
CREATE INDEX idx_categories_is_system ON categories(is_system);
