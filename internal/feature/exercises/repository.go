package exercises

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	constants "github.com/Sanchir01/fitnow/internal/models/contants"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type Repository struct {
	log *slog.Logger
	db  *pgxpool.Pool
}
type ExerciseDB struct {
	ID          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	VideoURL    string    `db:"video_url"`
}

func NewRepository(log *slog.Logger, db *pgxpool.Pool) *Repository {
	return &Repository{log: log, db: db}
}

func (r *Repository) GetAllExercises(ctx context.Context) ([]ExerciseDB, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	query, args, err := sq.Select("id", "title", "description", "video_url").From(constants.ExercisesTableName).ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	exercises := make([]ExerciseDB, 0)
	for rows.Next() {
		var item ExerciseDB
		if err := rows.Scan(&item.ID, &item.Title, &item.Description, &item.VideoURL); err != nil {
			return nil, err
		}
		exercises = append(exercises, item)
	}
	return exercises, nil
}

func (r *Repository) CreateExercise(ctx context.Context, title, desc string) (*ExerciseDB, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query, args, err := sq.
		Insert(constants.ExercisesTableName).
		Columns("title", "description").
		Values(title, desc).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id,title,description").
		ToSql()
	if err != nil {
		return nil, err
	}
	var exercise ExerciseDB
	if err := conn.QueryRow(ctx, query, args...).Scan(&exercise.ID, &exercise.Title, &exercise.Description); err != nil {
		return nil, err
	}
	return &exercise, nil
}

func (r *Repository) UpdateVideoUrl(ctx context.Context, url string, id uuid.UUID) error {
	query, args, err := sq.
		Update(constants.ExercisesTableName).
		Set("video_url", url).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no exercise was updated")
	}
	return nil
}
