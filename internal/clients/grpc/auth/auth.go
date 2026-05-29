package authgrpc

import (
	"context"
	"crypto/tls"
	"fmt"
	authv1 "github.com/Fitnow08/fitnow-proto/pkg/gen/go/v1/auth"
	"github.com/Sanchir01/fitnow/internal/config"
	"github.com/Sanchir01/fitnow/pkg/logger"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
)

type AuthClient struct {
	authv1.AuthServiceClient
	log *slog.Logger
}

func NewAuthClient(log *slog.Logger, cfg config.Client) (*AuthClient, error) {
	const op = "AuthClient.NewAuthClient"
	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}
	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.Aborted, codes.DeadlineExceeded, codes.NotFound),
		grpcretry.WithMax(uint(cfg.Retries)),
		grpcretry.WithPerRetryTimeout(cfg.Timeout),
	}

	creds := insecure.NewCredentials()
	if !cfg.Insecure {
		creds = credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})
	}

	cc, err := grpc.NewClient(
		cfg.Address,
		grpc.WithTransportCredentials(creds),
		grpc.WithChainUnaryInterceptor(
			grpcretry.UnaryClientInterceptor(retryOpts...),
			grpclog.UnaryClientInterceptor(logger.InterceptorLogger(log), logOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &AuthClient{
		AuthServiceClient: authv1.NewAuthServiceClient(cc),
		log:               log,
	}, nil
}

func (c *AuthClient) Login(ctx context.Context, email, password string) (*authv1.LoginResponse, error) {
	const op = "AuthClient.Login"
	resp, err := c.AuthServiceClient.Login(ctx, &authv1.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resp, nil
}

func (c *AuthClient) Register(ctx context.Context, email, password, name string) (*authv1.RegisterResponse, error) {
	const op = "AuthClient.Register"
	resp, err := c.AuthServiceClient.Register(ctx, &authv1.RegisterRequest{
		Email:    email,
		Password: password,
		Name:     name,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resp, nil
}

func (c *AuthClient) NewTokens(ctx context.Context, refreshToken string) (*authv1.NewTokensResponse, error) {
	const op = "AuthClient.NewTokens"
	resp, err := c.AuthServiceClient.NewTokens(ctx, &authv1.NewTokensRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resp, nil
}

func (c *AuthClient) VerifyAccount(ctx context.Context, email string, verifyCode int64) (*authv1.VerifyAccountResponse, error) {
	const op = "AuthClient.VerifyAccount"
	resp, err := c.AuthServiceClient.VerifyAccount(ctx, &authv1.VerifyAccountRequest{
		Email:      email,
		VerifyCode: verifyCode,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resp, nil
}

func (c *AuthClient) ResendVerifyCode(ctx context.Context, email string) (*authv1.ResendVerifyCodeResponse, error) {
	const op = "AuthClient.ResendVerifyCode"
	resp, err := c.AuthServiceClient.ResendVerifyCode(ctx, &authv1.ResendVerifyCodeRequest{
		Email: email,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resp, nil
}

func (c *AuthClient) ResetPassword(ctx context.Context, email string) (*authv1.ResetPasswordResponse, error) {
	const op = "AuthClient.ResetPassword"
	resp, err := c.AuthServiceClient.ResetPassword(ctx, &authv1.ResetPasswordRequest{
		Email: email,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resp, nil
}

func (c *AuthClient) ConfirmResetPassword(ctx context.Context, email, newPassword string, code int64) (*authv1.ConfirmResetPasswordResponse, error) {
	const op = "AuthClient.ConfirmResetPassword"
	resp, err := c.AuthServiceClient.ConfirmResetPassword(ctx, &authv1.ConfirmResetPasswordRequest{
		Email:       email,
		NewPassword: newPassword,
		Code:        code,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resp, nil
}
