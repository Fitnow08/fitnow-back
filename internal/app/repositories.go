package app

import (
	"github.com/Sanchir01/fitnow/internal/feature/auth"
	"log/slog"
)

type Repositories struct {
	AuthRepository *auth.Repository
}

func NewRepository(db *Database, l *slog.Logger) *Repositories {
	return &Repositories{
		AuthRepository: auth.NewRepository(l, db.PrimaryDB),
	}
}
