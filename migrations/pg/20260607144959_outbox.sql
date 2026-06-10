-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS events (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    event_type TEXT NOT NULL,
    payload TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'new' CHECK ( status IN ('new','done') ),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reserved_to TIMESTAMP DEFAULT  NULL
);
-- Частичный индекс для воркера-релейера: быстро выбирать неотправленные события.
CREATE INDEX IF NOT EXISTS idx_events_new ON events (created_at) WHERE status = 'new';

-- Автообновление updated_at при UPDATE (функция уже создаётся в миграции users).
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_events_updated_at BEFORE UPDATE
    ON events FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_events_updated_at ON events;
DROP INDEX IF EXISTS idx_events_new;
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
