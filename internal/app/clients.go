package app

import (
	authgrpc "github.com/Sanchir01/fitnow/internal/clients/grpc/auth"
	catalogrpc "github.com/Sanchir01/fitnow/internal/clients/grpc/catalog"
	"github.com/Sanchir01/fitnow/internal/config"
	"log/slog"
)

type Clients struct {
	AuthClient    *authgrpc.AuthClient
	ProgramClient *catalogrpc.ProgramClient
}

func NewClients(cfg *config.Config, l *slog.Logger) (*Clients, error) {
	authClient, err := authgrpc.NewAuthClient(l, cfg.Clients.Auth)
	if err != nil {
		l.Info("authClient", "err", err.Error())
		return nil, err
	}
	programClient, err := catalogrpc.NewProgramClient(l, cfg.Clients.CatalogClient)
	if err != nil {
		l.Info("programClient", "err", err.Error())
		return nil, err
	}
	return &Clients{
		AuthClient:    authClient,
		ProgramClient: programClient,
	}, nil
}
