-- +goose Up
-- +goose StatementBegin
CREATE TYPE subscription_status AS ENUM ('active', 'canceled', 'expired', 'pending');
CREATE TYPE billing_interval    AS ENUM ('month', 'year');

CREATE TABLE IF NOT EXISTS subscription_plans (
                                                  id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                                  code        TEXT NOT NULL UNIQUE,            -- 'free' | 'pro' | 'premium'
                                                  name        TEXT NOT NULL,
                                                  level       INT  NOT NULL DEFAULT 0,
                                                  price       BIGINT NOT NULL DEFAULT 0,
                                                  currency    TEXT NOT NULL DEFAULT 'RUB',
                                                  interval    billing_interval NOT NULL DEFAULT 'month',
                                                  is_active   BOOLEAN NOT NULL DEFAULT true,
                                                  created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                                  updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                                  version     BIGINT NOT NULL DEFAULT 1
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
