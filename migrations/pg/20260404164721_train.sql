-- +goose Up
-- +goose StatementBegin
CREATE TYPE train_type AS ENUM ('strength', 'cardio', 'stretching');

CREATE TYPE difficulty_level AS ENUM ('easy', 'medium', 'hard');

CREATE TABLE exercises (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    description TEXT
);
CREATE TABLE IF NOT EXISTS  trains (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY ,
    title TEXT NOT NULL,
    type train_type,
    duration BIGINT,
    is_public BOOLEAN DEFAULT true NOT NULL,
    difficulty difficulty_level,
    created_by UUID REFERENCES users(id),
    calories BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version BIGINT NOT NULL DEFAULT 0
);
CREATE TABLE IF NOT EXISTS train_exercises (
    id UUID DEFAULT uuid_generate_v4() NOT NULL,
    steps INT,
    sets INT,
    position INT NOT NULL,
    train_id UUID NOT NULL REFERENCES trains(id) ON DELETE CASCADE,
    exercises_id UUID NOT NULL REFERENCES exercises(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE user_trains (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    train_id UUID NOT NULL REFERENCES trains(id) ON DELETE CASCADE,
    added_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, train_id)
);

CREATE TRIGGER update_train_exercises_updated_at
    BEFORE UPDATE ON train_exercises
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_trains_updated_at
    BEFORE UPDATE ON trains
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_trains;
DROP TABLE IF EXISTS train_exercises;
DROP TABLE IF EXISTS exercises;
DROP TABLE IF EXISTS trains;
DROP TYPE IF EXISTS train_type;
DROP TYPE IF EXISTS difficulty_level;

-- +goose StatementEnd
