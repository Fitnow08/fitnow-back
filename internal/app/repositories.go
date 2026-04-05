package app

import (
	"github.com/Sanchir01/fitnow/internal/feature/auth"
	"github.com/Sanchir01/fitnow/internal/feature/train"
	"log/slog"
)

type Repositories struct {
	AuthRepository  *auth.Repository
	TrainRepository *train.Repository
}

func NewRepository(db *Database, l *slog.Logger) *Repositories {
	return &Repositories{
		AuthRepository:  auth.NewRepository(l, db.PrimaryDB),
		TrainRepository: train.NewRepository(db.PrimaryDB, l),
	}
}
