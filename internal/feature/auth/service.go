package auth

import (
	"context"
	"errors"
	authgrpc "github.com/Sanchir01/fitnow/internal/clients/grpc/auth"
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"github.com/google/uuid"
	"log/slog"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, email, title string, password []byte) (*UserDB, error)
	UserByEmail(ctx context.Context, email string) (*UserDB, error)
}
type Service struct {
	log            *slog.Logger
	authrepository AuthRepository
	authClient     *authgrpc.AuthClient
}

func NewService(log *slog.Logger, authrepository AuthRepository, authClient *authgrpc.AuthClient) *Service {
	return &Service{log: log, authrepository: authrepository, authClient: authClient}
}
func (s *Service) Register(ctx context.Context, req RegisterRequest) error {
	const op = "Auth.Service.Register"
	log := s.log.With("op", op)
	_, err := s.authClient.Register(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*domain.User, error) {
	const op = "Auth.Service.Login"
	log := s.log.With("op", op)
	data, err := s.authClient.Login(ctx, email, password)
	if err != nil {
		log.Error("failed to login", "error", err)
		return nil, err
	}
	dataid, err := uuid.Parse(data.Id)
	if err != nil {
		log.Error("failed to parse id", "error", err)
		return nil, err
	}
	return &domain.User{
		ID:           dataid,
		Email:        data.Email,
		Title:        data.Title,
		AccessToken:  data.RefreshToken,
		RefreshToken: data.RefreshToken,
	}, nil
}

func (s *Service) GenerateNewTokens(ctx context.Context, token string) (*Tokens, error) {
	const op = "Auth.Service.GenerateNewTokens"
	log := s.log.With("op", op)
	log.Info("Generating new tokens")

	tokens, err := s.authClient.NewTokens(ctx, token)
	if err != nil {
		log.Error("failed to generate tokens", "error", err)
		return nil, err
	}
	return &Tokens{
		RefreshToken: tokens.RefreshToken,
		AccessToken:  tokens.AccessToken,
	}, nil
}

func (s *Service) VerifyAccount(ctx context.Context, email string, code int64) (*domain.User, error) {
	const op = "Auth.Service.VerifyAccount"
	log := s.log.With("op", op)

	data, err := s.authClient.VerifyAccount(ctx, email, code)
	if err != nil {
		log.Error("failed to verify account", "error", err)
		return nil, err
	}
	uuidac, err := uuid.Parse(data.Id)
	if err != nil {
		log.Error("failed to parse id", "error", err)
		return nil, err
	}
	return &domain.User{
		uuidac,
		data.Email,
		data.Title,
		data.RefreshToken,
		data.AccessToken,
	}, nil
}

func (s *Service) ResendVerifyCode(ctx context.Context, email string) error {
	const op = "Auth.Service.ResendVerifyCode"
	log := s.log.With("op", op)
	_, err := s.authClient.ResendVerifyCode(ctx, email)
	if err != nil {
		log.Error("failed to resend verify code", "error", err)
		return err
	}
	return nil
}

func (s *Service) ResetPassword(ctx context.Context, email string) error {
	const op = "Auth.Service.ResetPassword"
	log := s.log.With("op", op)
	_, err := s.authClient.ResetPassword(ctx, email)
	if err != nil {
		log.Error("failed to reset password", "error", err)
		return err
	}
	return nil
}
func (s *Service) ConfirmResetPassword(ctx context.Context, email, newPassword string, code int64) error {
	const op = "Auth.Service.ConfirmResetPassword"
	log := s.log.With("op", op)
	ok, err := s.authClient.ConfirmResetPassword(ctx, email, newPassword, code)
	if err != nil {
		log.Error("failed to confirm reset password", "error", err)
		return err
	}
	if ok.Ok {
		return errors.New("failed to confirm reset password")
	}
	return nil
}
