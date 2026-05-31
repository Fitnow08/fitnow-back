-- +goose Up
-- +goose StatementBegin
ALTER TABLE programs ADD COLUMN IF NOT EXISTS image_path TEXT NOT NULL DEFAULT 'program.webp';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE programs DROP COLUMN IF EXISTS image_path;
-- +goose StatementEnd
