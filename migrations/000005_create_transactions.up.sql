CREATE TYPE transaction_type AS ENUM ('DEBIT', 'CREDIT');

CREATE TABLE transactions (
    id                      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id                 UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_id              UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    pluggy_transaction_id   VARCHAR(255) UNIQUE,
    description             VARCHAR(500) NOT NULL,
    amount                  NUMERIC(15, 2) NOT NULL,
    date                    DATE NOT NULL,
    type                    transaction_type NOT NULL,
    category_id             UUID REFERENCES categories(id) ON DELETE SET NULL,
    pluggy_category         VARCHAR(255),
    notes                   TEXT,
    tags                    TEXT[],
    is_recurring            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_account_id ON transactions(account_id);
CREATE INDEX idx_transactions_date ON transactions(date);
CREATE INDEX idx_transactions_category_id ON transactions(category_id);
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_pluggy_id ON transactions(pluggy_transaction_id);
CREATE INDEX idx_transactions_user_date ON transactions(user_id, date DESC);
