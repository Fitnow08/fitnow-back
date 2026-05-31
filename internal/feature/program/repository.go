package program

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
type Level string

const (
	Easy   Level = "easy"
	Medium Level = "medium"
	Hard   Level = "hard"
)

func NewRepository(log *slog.Logger, db *pgxpool.Pool) *Repository {
	return &Repository{log: log, db: db}
}

func (r *Repository) GetAllPrograms(ctx context.Context) ([]ProgramAndCountTrainDB, error) {
	query, args, err := sq.
		Select(
			"p.id", "p.title", "p.description", "p.weeks", "p.difficulty",
			"p.is_public", "p.version", "p.created_by", "p.created_at", "p.updated_at",
			"p.category_id", "p.image_path",
			"COUNT(pt.id) AS trains_count",
		).
		From(constants.ProgramTableName + " p").
		LeftJoin(constants.ProgramTrainsTableName + " pt ON pt.program_id = p.id").
		Where(sq.Eq{"p.is_public": true}).
		GroupBy("p.id").
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

	var programs []ProgramAndCountTrainDB
	for rows.Next() {
		var p ProgramAndCountTrainDB
		if err := rows.Scan(&p.ID, &p.Title, &p.Desc, &p.Weeks, &p.Difficult, &p.IsPublic, &p.Version, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt, &p.CategoryID, &p.ImagePath, &p.TrainsCount); err != nil {
			return nil, err
		}
		programs = append(programs, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return programs, nil
}
func (r *Repository) CreateProgram(ctx context.Context, title string, description string, weeks int, level Level, categoryID *uuid.UUID, user_id uuid.UUID) (*ProgramDB, error) {
	query, args, err := sq.
		Insert(constants.ProgramTableName).
		Columns("title", "description", "weeks", "difficulty", "category_id", "created_by").
		Values(title, description, weeks, level, categoryID, user_id).
		Suffix("RETURNING id, title, description, weeks, difficulty, is_public, category_id, image_path, created_by, created_at, updated_at, version").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var p ProgramDB
	if err := r.db.QueryRow(ctx, query, args...).Scan(
		&p.ID, &p.Title, &p.Desc, &p.Weeks, &p.Difficult, &p.IsPublic,
		&p.CategoryID, &p.ImagePath, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt, &p.Version,
	); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) UpdateProgramImagePath(ctx context.Context, programID uuid.UUID, imagePath string) error {
	query, args, err := sq.
		Update(constants.ProgramTableName).
		Set("image_path", imagePath).
		Where(sq.Eq{"id": programID}).
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
