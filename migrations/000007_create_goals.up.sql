CREATE TYPE goal_type AS ENUM ('SAVINGS', 'DEBT_PAYOFF', 'INVESTMENT', 'EMERGENCY_FUND', 'OTHER');

CREATE TABLE goals (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    target_amount   NUMERIC(15, 2) NOT NULL,
    current_amount  NUMERIC(15, 2) NOT NULL DEFAULT 0,
    deadline        DATE,
    type            goal_type NOT NULL DEFAULT 'SAVINGS',
    is_completed    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_goals_user_id ON goals(user_id);
