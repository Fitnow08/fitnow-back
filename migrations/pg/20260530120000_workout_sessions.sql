-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS workout_sessions (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id          UUID NOT NULL REFERENCES users(id)  ON DELETE CASCADE,
    train_id         UUID NOT NULL REFERENCES trains(id) ON DELETE CASCADE,
    program_train_id UUID REFERENCES program_trains(id), -- NULL = разовая тренировка
    duration         BIGINT NOT NULL CHECK (duration >= 0),
    calories         BIGINT NOT NULL CHECK (calories >= 0),
    completed_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- выборка прогресса за период
CREATE INDEX idx_workout_sessions_user_completed
    ON workout_sessions (user_id, completed_at);

-- дедуп: один трейн = один зачёт в день (без учёта программы)
CREATE UNIQUE INDEX uniq_session_per_day
    ON workout_sessions (user_id, train_id, (completed_at::date));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS workout_sessions;
DROP INDEX IF EXISTS uniq_session_per_day;
DROP INDEX IF EXISTS idx_workout_sessions_user_completed;
-- +goose StatementEnd
