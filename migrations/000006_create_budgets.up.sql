CREATE TYPE budget_period AS ENUM ('MONTHLY', 'YEARLY');

CREATE TABLE budgets (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    name        VARCHAR(255) NOT NULL,
    amount      NUMERIC(15, 2) NOT NULL,
    period      budget_period NOT NULL DEFAULT 'MONTHLY',
    start_date  DATE NOT NULL,
    end_date    DATE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_budgets_user_id ON budgets(user_id);
CREATE INDEX idx_budgets_category_id ON budgets(category_id);
