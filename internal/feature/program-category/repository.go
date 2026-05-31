package programcategory

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	constants "github.com/Sanchir01/fitnow/internal/models/contants"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type Repository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewRepository(db *pgxpool.Pool, log *slog.Logger) *Repository {
	return &Repository{db: db, log: log}
}

func (r *Repository) GetAllProgramCategory(ctx context.Context) ([]ProgramCategoryDB, error) {
	query, args, err := sq.Select("id", "title", "created_at", "updated_at", "version").
		From(constants.ProgramCategoryTableName).
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
	categories := make([]ProgramCategoryDB, 0)
	for rows.Next() {
		var category ProgramCategoryDB
		if err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt, &category.Version); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *Repository) CreateProgramCategory(ctx context.Context, title string) (*ProgramCategoryDB, error) {
	query, args, err := sq.Insert(constants.ProgramCategoryTableName).
		Columns("title").
		Values(title).
		Suffix("RETURNING id, title, created_at, updated_at, version").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	var category ProgramCategoryDB
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt, &category.Version,
	)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *Repository) UpdateProgramCategory(ctx context.Context, id uuid.UUID, title string) (*ProgramCategoryDB, error) {
	query, args, err := sq.
		Update(constants.ProgramCategoryTableName).
		Set("title", title).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, title, created_at, updated_at, version").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	var category ProgramCategoryDB
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt, &category.Version,
	)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *Repository) DeleteProgramCategory(ctx context.Context, id uuid.UUID) error {
	query, args, err := sq.Delete(constants.ProgramCategoryTableName).
		Where(sq.Eq{"id": id}).
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
