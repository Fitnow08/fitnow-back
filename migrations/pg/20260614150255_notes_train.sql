-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS session_notes (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID NOT NULL REFERENCES users(id)            ON DELETE CASCADE,
    session_id UUID NOT NULL REFERENCES workout_sessions(id) ON DELETE CASCADE,
    note       TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- заметки конкретного занятия, новые сверху
CREATE INDEX IF NOT EXISTS idx_session_notes_session
    ON session_notes (session_id, created_at DESC);

-- все заметки пользователя
CREATE INDEX IF NOT EXISTS idx_session_notes_user
    ON session_notes (user_id, created_at DESC);

CREATE TRIGGER update_session_notes_updated_at
    BEFORE UPDATE ON session_notes
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_session_notes_updated_at ON session_notes;
DROP INDEX IF EXISTS idx_session_notes_session;
DROP INDEX IF EXISTS idx_session_notes_user;
DROP TABLE IF EXISTS session_notes;
-- +goose StatementEnd
