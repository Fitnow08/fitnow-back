-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS programs (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title       TEXT NOT NULL,
    description TEXT,
    weeks       INT  NOT NULL CHECK (weeks > 0),  -- длительность программы в неделях
    difficulty  difficulty_level,                 -- переиспользуем существующий ENUM
    is_public   BOOLEAN NOT NULL DEFAULT true,
    created_by  UUID REFERENCES users(id),
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version     BIGINT NOT NULL DEFAULT 0
);

CREATE TRIGGER update_programs_updated_at
    BEFORE UPDATE ON programs
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


CREATE TABLE IF NOT EXISTS program_trains (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    program_id  UUID NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    train_id    UUID NOT NULL REFERENCES trains(id)   ON DELETE CASCADE,
    week_number INT NOT NULL CHECK (week_number >= 1),           -- 1..programs.weeks
    day_of_week INT NOT NULL CHECK (day_of_week BETWEEN 1 AND 7), -- 1=Пн .. 7=Вс
    position    INT NOT NULL,                                    -- порядок в рамках дня
    UNIQUE (program_id, week_number, day_of_week, position)
);

CREATE INDEX idx_program_trains_program ON program_trains (program_id, week_number);

CREATE TABLE IF NOT EXISTS user_programs (
    user_id    UUID NOT NULL REFERENCES users(id)    ON DELETE CASCADE,
    program_id UUID NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_active  BOOLEAN NOT NULL DEFAULT true,
    PRIMARY KEY (user_id, program_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_programs;
DROP TABLE IF EXISTS program_trains;
DROP TRIGGER IF EXISTS update_programs_updated_at ON programs;
DROP TABLE IF EXISTS programs;
DROP INDEX IF EXISTS idx_program_trains_program;
-- +goose StatementEnd
