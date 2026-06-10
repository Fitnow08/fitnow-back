package app

import (
	"github.com/Sanchir01/fitnow/internal/feature/comment"
	"github.com/Sanchir01/fitnow/internal/feature/events"
	"github.com/Sanchir01/fitnow/internal/feature/exercises"
	programcategory "github.com/Sanchir01/fitnow/internal/feature/program-category"
	"github.com/Sanchir01/fitnow/internal/feature/rating"
	"github.com/Sanchir01/fitnow/internal/feature/train"
	traincategory "github.com/Sanchir01/fitnow/internal/feature/train-category"
	"log/slog"
)

type Repositories struct {
	TrainRepository           *train.Repository
	ExercisesRepository       *exercises.Repository
	TrainCategoryRepository   *traincategory.Repository
	RatingRepository          *rating.Repository
	CommentRepository         *comment.Repository
	ProgramCategoryRepository *programcategory.Repository
	EventRepository           *events.Repository
}

func NewRepository(db *Database, l *slog.Logger) *Repositories {
	return &Repositories{
		TrainRepository:           train.NewRepository(db.PrimaryDB, l),
		ExercisesRepository:       exercises.NewRepository(l, db.PrimaryDB),
		TrainCategoryRepository:   traincategory.NewRepository(db.PrimaryDB, l),
		CommentRepository:         comment.NewRepository(l, db.PrimaryDB),
		RatingRepository:          rating.NewRepository(l, db.PrimaryDB),
		ProgramCategoryRepository: programcategory.NewRepository(db.PrimaryDB, l),
		EventRepository:           events.NewRepository(l, db.PrimaryDB),
	}
}
