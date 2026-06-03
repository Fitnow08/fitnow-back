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

func (r *Repository) GetAllPrograms(ctx context.Context, categoryId uuid.UUID, search string) ([]ProgramAndCountTrainDB, error) {
	builder := sq.
		Select(
			"p.id", "p.title", "p.description", "p.weeks", "p.difficulty",
			"p.is_public", "p.version", "p.created_by", "p.created_at", "p.updated_at",
			"p.category_id", "p.image_path",
			"COUNT(pt.id) AS trains_count",
		).
		From(constants.ProgramTableName + " p").
		LeftJoin(constants.ProgramTrainsTableName + " pt ON pt.program_id = p.id").
		Where(sq.Eq{"p.is_public": true}).
		GroupBy("p.id")
	if categoryId != uuid.Nil {
		builder = builder.Where(sq.Eq{
			"p.category_id": categoryId,
		})
	}
	if search != "" {
		builder = builder.Where(sq.Or{
			sq.ILike{"p.title": "%" + search + "%"},
			sq.ILike{"p.description": "%" + search + "%"},
		})
	}
	query, args, err := builder.
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

func (r *Repository) AddProgramTrains(ctx context.Context, programId uuid.UUID, trains []ProgramTrainInput) error {
	if len(trains) == 0 {
		return nil
	}

	b := sq.
		Insert(constants.ProgramTrainsTableName).
		Columns("program_id", "train_id", "week_number", "day_of_week", "position").
		PlaceholderFormat(sq.Dollar)
	for _, train := range trains {
		b = b.Values(programId, train.TrainId, train.WeekNumber, train.DayOfWeek, train.Position)
	}
	query, args, err := b.ToSql()

	if err != nil {
		return err

	}
	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetProgramTrains(ctx context.Context, programId uuid.UUID) ([]ProgramTrainDB, error) {
	query, args, err := sq.
		Select(
			"t.id", "t.title", "t.type", "t.duration", "t.is_public",
			"t.difficulty", "t.calories", "t.created_by", "t.created_at",
			"t.image_path", "t.category_id",
			"pt.week_number", "pt.day_of_week", "pt.position",
		).
		From(constants.ProgramTrainsTableName+" pt").
		Join(constants.TrainTableName+" t ON t.id = pt.train_id").
		Where(sq.Eq{"pt.program_id": programId}).
		OrderBy("pt.week_number", "pt.day_of_week", "pt.position").
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

	var trains []ProgramTrainDB
	for rows.Next() {
		var t ProgramTrainDB
		if err := rows.Scan(
			&t.ID, &t.Title, &t.Type, &t.Duration, &t.IsPublic,
			&t.Difficulty, &t.Calories, &t.CreatedBy, &t.CreatedAt,
			&t.ImageURL, &t.CategoryId,
			&t.WeekNumber, &t.DayOfWeek, &t.Position,
		); err != nil {
			return nil, err
		}
		trains = append(trains, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return trains, nil
}

func (r *Repository) GetProgramById(ctx context.Context, programID uuid.UUID) (*ProgramDB, error) {
	query, args, err := sq.
		Select(
			"id", "title", "description", "weeks", "difficulty", "is_public",
			"category_id", "image_path", "created_by", "created_at", "updated_at", "version",
		).
		From(constants.ProgramTableName).
		Where(sq.Eq{"id": programID}).
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
