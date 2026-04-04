package auth

import (
	"context"
	"errors"
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, email, title string, password []byte) (*UserDB, error)
	UserByEmail(ctx context.Context, email string) (*UserDB, error)
}
type Service struct {
	log            *slog.Logger
	authrepository AuthRepository
}

func NewService(log *slog.Logger, authrepository AuthRepository) *Service {
	return &Service{log: log, authrepository: authrepository}
}
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*domain.User, error) {
	const op = "Auth.Service.Register"
	log := s.log.With("op", op)
	_, err := s.authrepository.UserByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.New("user already exists")
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err // только реальные ошибки                                                                                                                                                                                    }
	}

	hashpass, err := GeneratePasswordHash(req.Password)
	if err != nil {
		log.Error("failed to generate password hash", "error", err)
		return nil, err
	}
	userdb, err := s.authrepository.CreateUser(ctx, req.Email, req.Name, hashpass)
	if err != nil {
		log.Error("failed to create user", "error", err)
		return nil, err
	}

	return &domain.User{
		ID:    userdb.ID,
		Email: userdb.Email,
		Title: userdb.Title,
	}, nil
}
