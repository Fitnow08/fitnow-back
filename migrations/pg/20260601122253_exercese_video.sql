-- +goose Up
-- +goose StatementBegin
ALTER TABLE exercises ADD COLUMN IF NOT EXISTS  video_url TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE exercises DROP COLUMN IF EXISTS video_url;
-- +goose StatementEnd
