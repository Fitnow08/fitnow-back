package train

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

type Repository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

type TrainDB struct {
	ID         uuid.UUID `db:"id"`
	Title      string    `db:"title"`
	Type       string    `db:"type"`
	Duration   int64     `db:"duration"`
	IsPublic   bool      `db:"is_public"`
	Difficulty string    `db:"difficulty"`
	CreatedBy  uuid.UUID `db:"created_by"`
	CreatedAt  time.Time `db:"created_at"`
	Calories   int64     `db:"calories"`
	UpdatedAt  time.Time `db:"updated_at"`
	Version    int64     `db:"version"`
}

type ExerciseDB struct {
	ID          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
}

func NewRepository(db *pgxpool.Pool, log *slog.Logger) *Repository {
	return &Repository{db: db, log: log}
}

func (r *Repository) GetAllPublicTrains(ctx context.Context) ([]*TrainDB, error) {
	query, args, err := sq.
		Select("id", "title", "type", "duration", "is_public", "difficulty", "calories", "created_by", "created_at", "version").
		From("trains").
		Where(sq.Eq{"is_public": true}).
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

	var trains []*TrainDB
	for rows.Next() {
		var t TrainDB
		if err := rows.Scan(&t.ID, &t.Title, &t.Type, &t.Duration, &t.IsPublic, &t.Difficulty, &t.Calories, &t.CreatedBy, &t.CreatedAt, &t.Version); err != nil {
			return nil, err
		}
		trains = append(trains, &t)
	}
	return trains, nil
}

func (r *Repository) GetTrainByID(ctx context.Context, id uuid.UUID) (*TrainDB, error) {
	query, args, err := sq.
		Select("id", "title", "type", "duration", "is_public", "difficulty", "calories", "created_by", "created_at", "version").
		From("trains").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	var t TrainDB
	if err := r.db.QueryRow(ctx, query, args...).Scan(&t.ID, &t.Title, &t.Type, &t.Duration, &t.IsPublic, &t.Difficulty, &t.Calories, &t.CreatedBy, &t.CreatedAt, &t.Version); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repository) CreateTrain(ctx context.Context, req CreateTrainRequest, userID uuid.UUID) (*TrainDB, error) {
	query, args, err := sq.
		Insert("trains").
		Columns("title", "type", "duration", "is_public", "difficulty", "calories", "created_by").
		Values(req.Title, req.Type, req.Duration, req.IsPublic, req.Difficulty, req.Calories, userID).
		Suffix("RETURNING id, title, type, duration, is_public, difficulty, calories, created_by, created_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	var t TrainDB
	if err := r.db.QueryRow(ctx, query, args...).Scan(&t.ID, &t.Title, &t.Type, &t.Duration, &t.IsPublic, &t.Difficulty, &t.Calories, &t.CreatedBy, &t.CreatedAt); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repository) UpdateTrain(ctx context.Context, id uuid.UUID, req UpdateTrainRequest, userID uuid.UUID) (*TrainDB, error) {
	builder := sq.Update("trains").Where(sq.Eq{"id": id, "created_by": userID}).PlaceholderFormat(sq.Dollar)
	if req.Title != "" {
		builder = builder.Set("title", req.Title)
	}
	if req.Type != "" {
		builder = builder.Set("type", req.Type)
	}
	if req.Duration > 0 {
		builder = builder.Set("duration", req.Duration)
	}
	if req.IsPublic != nil {
		builder = builder.Set("is_public", *req.IsPublic)
	}
	if req.Difficulty != "" {
		builder = builder.Set("difficulty", req.Difficulty)
	}
	if req.Calories > 0 {
		builder = builder.Set("calories", req.Calories)
	}
	query, args, err := builder.
		Suffix("RETURNING id, title, type, duration, is_public, difficulty, calories, created_by, created_at").
		ToSql()
	if err != nil {
		return nil, err
	}
	var t TrainDB
	if err := r.db.QueryRow(ctx, query, args...).Scan(&t.ID, &t.Title, &t.Type, &t.Duration, &t.IsPublic, &t.Difficulty, &t.Calories, &t.CreatedBy, &t.CreatedAt); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repository) DeleteTrain(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query, args, err := sq.
		Delete("trains").
		Where(sq.Eq{"id": id, "created_by": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, query, args...)
	return err
}

func (r *Repository) GetUserTrains(ctx context.Context, userID uuid.UUID) ([]*TrainDB, error) {
	query, args, err := sq.
		Select("t.id", "t.title", "t.type", "t.duration", "t.is_public", "t.difficulty", "t.calories", "t.created_by", "t.created_at").
		From("trains t").
		Join("user_trains ut ON ut.train_id = t.id").
		Where(sq.Eq{"ut.user_id": userID}).
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

	var trains []*TrainDB
	for rows.Next() {
		var t TrainDB
		if err := rows.Scan(&t.ID, &t.Title, &t.Type, &t.Duration, &t.IsPublic, &t.Difficulty, &t.Calories, &t.CreatedBy, &t.CreatedAt); err != nil {
			return nil, err
		}
		trains = append(trains, &t)
	}
	return trains, nil
}

func (r *Repository) AddUserTrain(ctx context.Context, userID, trainID uuid.UUID) error {
	query, args, err := sq.
		Insert("user_trains").
		Columns("user_id", "train_id").
		Values(userID, trainID).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, query, args...)
	return err
}

func (r *Repository) RemoveUserTrain(ctx context.Context, userID, trainID uuid.UUID) error {
	query, args, err := sq.
		Delete("user_trains").
		Where(sq.Eq{"user_id": userID, "train_id": trainID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, query, args...)
	return err
}

func (r *Repository) GetAllExercises(ctx context.Context) ([]*ExerciseDB, error) {
	query, args, err := sq.
		Select("id", "title", "description").
		From("exercises").
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

	var exercises []*ExerciseDB
	for rows.Next() {
		var e ExerciseDB
		if err := rows.Scan(&e.ID, &e.Title, &e.Description); err != nil {
			return nil, err
		}
		exercises = append(exercises, &e)
	}
	return exercises, nil
}

func (r *Repository) CreateExercise(ctx context.Context, req CreateExerciseRequest) (*ExerciseDB, error) {
	query, args, err := sq.
		Insert("exercises").
		Columns("title", "description").
		Values(req.Title, req.Description).
		Suffix("RETURNING id, title, description").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	var e ExerciseDB
	if err := r.db.QueryRow(ctx, query, args...).Scan(&e.ID, &e.Title, &e.Description); err != nil {
		return nil, err
	}
	return &e, nil
}
