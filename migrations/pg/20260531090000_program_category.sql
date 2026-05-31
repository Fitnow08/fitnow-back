-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS program_category (
    id   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version BIGINT NOT NULL DEFAULT 1
);

ALTER TABLE programs
    ADD COLUMN category_id UUID REFERENCES program_category(id) ON DELETE SET NULL;

-- ускоряем фильтрацию программ по категории
CREATE INDEX idx_programs_category ON programs (category_id);

CREATE TRIGGER update_program_category_updated_at
    BEFORE UPDATE ON program_category
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_programs_category;
ALTER TABLE programs DROP COLUMN IF EXISTS category_id;
DROP TRIGGER IF EXISTS update_program_category_updated_at ON program_category;
DROP TABLE IF EXISTS program_category;
-- +goose StatementEnd
