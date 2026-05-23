package comment

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	constants "github.com/Sanchir01/fitnow/internal/models/contants"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

type Repository struct {
	log *slog.Logger
	db  *pgxpool.Pool
}

func NewRepository(log *slog.Logger, db *pgxpool.Pool) *Repository {
	return &Repository{log: log, db: db}
}

func (r *Repository) CreateComment(ctx context.Context, comment string, train_id, user_id uuid.UUID, parentid *uuid.UUID) error {
	query, args, err := sq.
		Insert(constants.CommentsTableName).
		Columns("train_id", "user_id", "parent_id", "content").
		Values(train_id, user_id, parentid, comment).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetTrainComments(ctx context.Context, trainID uuid.UUID) ([]CommentDB, error) {
	query, args, err := sq.
		Select(
			"id",
			"train_id",
			"user_id",
			"parent_id",
			"content",
			"created_at",
			"updated_at",
			"is_deleted",
			"deleted_at",
			"version",
		).
		From(constants.CommentsTableName).
		Where(sq.Eq{
			"train_id": trainID,
		}).
		OrderBy("created_at ASC").
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

	comments := make([]CommentDB, 0)

	for rows.Next() {
		var row CommentDB

		if err := rows.Scan(
			&row.ID,
			&row.TrainID,
			&row.UserId,
			&row.ParentID,
			&row.Comment,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.IsDeleted,
			&row.DeletedAt,
			&row.Version,
		); err != nil {
			return nil, err
		}

		comments = append(comments, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *Repository) UpdateComment(ctx context.Context, comment string, commentID uuid.UUID) error {
	query, args, err := sq.Update(constants.CommentsTableName).
		Set("content", comment).
		Where(sq.Eq{"id": commentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteComment(ctx context.Context, commentID uuid.UUID) error {
	queery, args, err := sq.Update(constants.CommentsTableName).
		Set("is_deleted", true).
		Set("deleted_at", time.Now()).
		Where(sq.Eq{"id": commentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, queery, args...)
	if err != nil {
		return err
	}
	return nil
}
