package app

import (
	"github.com/Sanchir01/fitnow/internal/feature/auth"
	"log/slog"
)

type Handlers struct {
	AuthHandler *auth.Handler
}

func NewHandlers(l *slog.Logger, srv *Services) *Handlers {
	return &Handlers{
		AuthHandler: auth.NewHandler(l, srv.AuthService),
	}
}
