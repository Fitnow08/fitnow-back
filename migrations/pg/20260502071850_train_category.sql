-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS train_category (
    id   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL UNIQUE ,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version BIGINT  NOT NULL DEFAULT  1
);

ALTER TABLE trains
    ADD COLUMN category_id UUID REFERENCES train_category(id) ON DELETE SET NULL;

CREATE TRIGGER update_train_category_updated_at
    BEFORE UPDATE ON train_category
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE trains DROP COLUMN IF EXISTS category_id;
DROP TRIGGER IF EXISTS update_train_category_updated_at ON train_category;
DROP TABLE IF EXISTS train_category;
-- +goose StatementEnd
