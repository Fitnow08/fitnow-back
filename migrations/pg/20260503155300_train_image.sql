-- +goose Up
-- +goose StatementBegin
ALTER TABLE trains ADD COLUMN  IF NOT EXISTS image_path TEXT NOT NULL DEFAULT 'train.webp';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE trains DROP COLUMN IF EXISTS image_path;
-- +goose StatementEnd
