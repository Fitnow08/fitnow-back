package app

import (
	"github.com/Sanchir01/fitnow/internal/feature/auth"
	"github.com/Sanchir01/fitnow/internal/feature/comment"
	"github.com/Sanchir01/fitnow/internal/feature/exercises"
	"github.com/Sanchir01/fitnow/internal/feature/rating"
	"github.com/Sanchir01/fitnow/internal/feature/train"
	traincategory "github.com/Sanchir01/fitnow/internal/feature/train-category"
	"log/slog"
)

type Repositories struct {
	AuthRepository          *auth.Repository
	TrainRepository         *train.Repository
	ExercisesRepository     *exercises.Repository
	TrainCategoryRepository *traincategory.Repository
	RatingRepository        *rating.Repository
	CommentRepository       *comment.Repository
}

func NewRepository(db *Database, l *slog.Logger) *Repositories {
	return &Repositories{
		AuthRepository:          auth.NewRepository(l, db.PrimaryDB),
		TrainRepository:         train.NewRepository(db.PrimaryDB, l),
		ExercisesRepository:     exercises.NewRepository(l, db.PrimaryDB),
		TrainCategoryRepository: traincategory.NewRepository(db.PrimaryDB, l),
		CommentRepository:       comment.NewRepository(l, db.PrimaryDB),
		RatingRepository:        rating.NewRepository(l, db.PrimaryDB),
	}
}
