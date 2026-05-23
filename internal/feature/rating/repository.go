package rating

import (
	"context"
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

func NewRepository(log *slog.Logger, db *pgxpool.Pool) *Repository {
	return &Repository{log: log, db: db}
}

func (r *Repository) GetAllRatings(ctx context.Context) ([]RatingDB, error) {
	query, args, err := sq.
		Select("user_id", "train_id", "rating", "created_at", "updated_at").
		From(constants.RatingTableName).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ratings := make([]RatingDB, 0)
	for rows.Next() {
		var rating RatingDB

		if err := rows.Scan(
			&rating.UserId,
			&rating.TrainId,
			&rating.Rating,
			&rating.CreatedAt,
			&rating.UpdatedAt,
		); err != nil {
			return nil, err
		}

		ratings = append(ratings, rating)
	}
	return ratings, nil
}
func (r *Repository) CreateTrainRating(ctx context.Context, userid, trainid uuid.UUID, rating int) error {
	query, args, err := sq.Insert(constants.RatingTableName).
		Columns("user_id", "train_id", "rating").
		Values(userid, trainid, rating).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	if _, err = r.db.Exec(ctx, query, args...); err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateTrainRating(ctx context.Context, userid, trainid uuid.UUID, rating int) error {
	query, args, err := sq.Insert(constants.RatingTableName).
		Columns("user_id", "train_id", "rating").
		Values(userid, trainid, rating).
		Suffix("ON CONFLICT (user_id, train_id) DO UPDATE SET rating = EXCLUDED.rating, updated_at = CURRENT_TIMESTAMP").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	if _, err = r.db.Exec(ctx, query, args...); err != nil {
		return err
	}
	return nil
}
