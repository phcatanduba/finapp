CREATE TYPE account_type AS ENUM ('BANK', 'CREDIT', 'INVESTMENT', 'LOAN', 'OTHER');

CREATE TABLE accounts (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_id             UUID NOT NULL REFERENCES pluggy_items(id) ON DELETE CASCADE,
    pluggy_account_id   VARCHAR(255) NOT NULL UNIQUE,
    name                VARCHAR(255) NOT NULL,
    type                account_type NOT NULL DEFAULT 'BANK',
    subtype             VARCHAR(100),
    balance             NUMERIC(15, 2) NOT NULL DEFAULT 0,
    credit_limit        NUMERIC(15, 2),
    currency_code       VARCHAR(3) NOT NULL DEFAULT 'BRL',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_accounts_user_id ON accounts(user_id);
CREATE INDEX idx_accounts_item_id ON accounts(item_id);
CREATE INDEX idx_accounts_pluggy_account_id ON accounts(pluggy_account_id);
