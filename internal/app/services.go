package app

import (
	"github.com/Sanchir01/fitnow/internal/feature/auth"
	"github.com/Sanchir01/fitnow/internal/feature/train"
	"log/slog"
)

type Services struct {
	AuthService  *auth.Service
	TrainService *train.Service
}

func NewServices(repo *Repositories, l *slog.Logger) *Services {
	return &Services{
		AuthService:  auth.NewService(l, repo.AuthRepository),
		TrainService: train.NewService(l, repo.TrainRepository),
	}
}
