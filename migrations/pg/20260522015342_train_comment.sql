-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS comments (
                                        id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                        train_id   UUID NOT NULL REFERENCES trains(id) ON DELETE CASCADE,
                                        user_id    UUID NOT NULL REFERENCES users(id)  ON DELETE CASCADE,
                                        parent_id  UUID REFERENCES comments(id) ON DELETE CASCADE,
                                        content    TEXT NOT NULL CHECK (char_length(content) BETWEEN 1 AND 2000),
                                        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        is_deleted BOOLEAN DEFAULT false,
                                        deleted_at TIMESTAMP,
                                        version    BIGINT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_comments_train_id   ON comments(train_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_comments_user_id    ON comments(user_id);
CREATE INDEX IF NOT EXISTS idx_comments_parent_id  ON comments(parent_id);


CREATE TRIGGER update_comments_updated_at
    BEFORE UPDATE ON comments
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS comments;
DROP INDEX  IF EXISTS idx_comments_train_id;
DROP INDEX  IF EXISTS idx_comments_user_id;
DROP INDEX  IF EXISTS idx_comments_parent_id;
-- +goose StatementEnd
