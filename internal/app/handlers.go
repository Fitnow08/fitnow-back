package app

import (
	"github.com/Sanchir01/fitnow/internal/feature/auth"
	"github.com/Sanchir01/fitnow/internal/feature/train"
	"log/slog"
)

type Handlers struct {
	AuthHandler  *auth.Handler
	TrainHandler *train.Handler
}

func NewHandlers(l *slog.Logger, srv *Services) *Handlers {
	return &Handlers{
		AuthHandler:  auth.NewHandler(l, srv.AuthService),
		TrainHandler: train.NewHandler(l, srv.TrainService),
	}
}
