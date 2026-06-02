package user

import (
	"context"
	authv1 "github.com/Fitnow08/fitnow-proto/pkg/gen/go/v1/auth"
	authgrpc "github.com/Sanchir01/fitnow/internal/clients/grpc/auth"
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"github.com/google/uuid"
	"log/slog"
)

type Service struct {
	log        *slog.Logger
	authClient *authgrpc.AuthClient
}

func NewService(log *slog.Logger, authClient *authgrpc.AuthClient) *Service {
	return &Service{log: log, authClient: authClient}
}

func (s *Service) GetAllUsers(ctx context.Context) ([]domain.OnlyUser, error) {
	const op = "AuthService.Service.GetAllUsers"
	log := s.log.With("op", op)
	users, err := s.authClient.GetAllUsers(ctx)
	if err != nil {
		log.Error("failed to get all users", "error", err)
		return nil, err
	}
	domainUsers := make([]domain.OnlyUser, 0, len(users.Users))
	for _, user := range users.Users {
		useruuid, err := uuid.Parse(user.Id)
		if err != nil {
			log.Error("failed to parse id", "error", err)
		}
		domainUsers = append(domainUsers, domain.OnlyUser{
			ID:    useruuid,
			Role:  user.Role,
			Title: user.Title,
			Email: user.Email,
		})
	}
	return domainUsers, nil
}
func (s *Service) GetUserById(ctx context.Context, id uuid.UUID) (*domain.OnlyUser, error) {
	const op = "AuthService.Service.GetUserByEmail"
	log := s.log.With("op", op)
	user, err := s.authClient.GetUserById(ctx, &authv1.GetUserByIdRequest{Id: id.String()})
	if err != nil {
		log.Error("failed to get user by id", "error", err)
		return nil, err
	}
	useruuid, err := uuid.Parse(user.User.Id)
	if err != nil {
		log.Error("failed to parse id", "error", err)
		return nil, err
	}
	return &domain.OnlyUser{
		ID:    useruuid,
		Email: user.User.Email,
		Title: user.User.Title,
		Role:  user.User.Role,
	}, err
}
