-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS train_ratings (
                                             user_id    UUID NOT NULL REFERENCES users(id)  ON DELETE CASCADE,
                                             train_id   UUID NOT NULL REFERENCES trains(id) ON DELETE CASCADE,
                                             rating     SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
                                             created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                             updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                             PRIMARY KEY (user_id, train_id)
);

CREATE INDEX IF NOT EXISTS idx_train_ratings_train_id ON train_ratings(train_id);

CREATE TRIGGER update_train_ratings_updated_at
    BEFORE UPDATE ON train_ratings
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS train_ratings;
DROP INDEX IF EXISTS idx_train_ratings_train_id
-- +goose StatementEnd
