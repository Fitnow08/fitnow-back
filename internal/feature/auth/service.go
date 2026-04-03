package auth

import (
	"context"
	"errors"
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, email, title string, password []byte, tx pgx.Tx) (*UserDB, error)
	UserByEmail(ctx context.Context, email string) (*UserDB, error)
}
type Service struct {
	log            *slog.Logger
	authrepository AuthRepository
}

func NewServiceAuth(log *slog.Logger, authrepository AuthRepository) *Service {
	return &Service{log: log, authrepository: authrepository}
}
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*domain.User, error) {
	user, err := s.authrepository.UserByEmail(ctx, req.Email)
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if user.ID != uuid.Nil {
		return nil, errors.New("user already exists")
	}
	
	return nil, nil
}
