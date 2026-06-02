-- +goose Up
-- +goose StatementBegin
ALTER TABLE exercises ADD COLUMN  IF NOT EXISTS  difficulty  difficulty_level NOT NULL DEFAULT 'easy';
ALTER TABLE train_exercises
    ADD CONSTRAINT train_exercises_train_id_exercise_id_unique
        UNIQUE (train_id, exercises_id);
ALTER TABLE program_trains
    ADD CONSTRAINT  program_trains_train_id_id_unique
        UNIQUE (train_id,program_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE train_exercises
    DROP CONSTRAINT IF EXISTS train_exercises_train_id_exercise_id_unique;

ALTER TABLE program_trains
    DROP CONSTRAINT IF EXISTS program_trains_train_id_id_unique;

ALTER TABLE exercises
    DROP COLUMN IF EXISTS difficulty;
-- +goose StatementEnd
