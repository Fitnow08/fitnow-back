package app

import (
	"github.com/Sanchir01/fitnow/internal/feature/auth"
	"github.com/Sanchir01/fitnow/internal/feature/exercises"
	"github.com/Sanchir01/fitnow/internal/feature/train"
	traincategory "github.com/Sanchir01/fitnow/internal/feature/train-category"
	"log/slog"
)

type Handlers struct {
	AuthHandler          *auth.Handler
	TrainHandler         *train.Handler
	ExercisesHandler     *exercises.Handler
	TrainCategoryHandler *traincategory.Handler
}

func NewHandlers(l *slog.Logger, srv *Services) *Handlers {
	return &Handlers{
		AuthHandler:          auth.NewHandler(l, srv.AuthService),
		TrainHandler:         train.NewHandler(l, srv.TrainService),
		ExercisesHandler:     exercises.NewHandler(l, srv.ExercisesService),
		TrainCategoryHandler: traincategory.NewHandler(l, srv.TrainCategoryService),
	}
}
