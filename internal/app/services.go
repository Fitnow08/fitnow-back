package app

import (
	"github.com/Sanchir01/fitnow/internal/feature/auth"
	"log/slog"
)

type Services struct {
	AuthService *auth.Service
}

func NewServices(repo *Repositories, l *slog.Logger) *Services {
	return &Services{
		AuthService: auth.NewService(l, repo.AuthRepository),
	}
}
