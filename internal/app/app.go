package app

import (
	"context"
	"fmt"
	authgrpc "github.com/Sanchir01/fitnow/internal/clients/grpc/auth"
	"github.com/Sanchir01/fitnow/internal/config"
	httpserver "github.com/Sanchir01/fitnow/internal/servers/http"
	"github.com/Sanchir01/fitnow/pkg/logger"
	"log/slog"
)

type App struct {
	Handlers     *Handlers
	Cfg          *config.Config
	HTTPServer   *httpserver.Server
	Log          *slog.Logger
	CancelLogger func()
}

func NewApp(ctx context.Context) (*App, error) {

	cfg := config.InitConfig()
	l, cancelogger := logger.SetupLogger(ctx, cfg.Env, fmt.Sprintf("%s:%s", "", ""))
	databases, err := NewDataBases(cfg, l)
	if err != nil {
		l.Info("databases", "err", err.Error())
		return nil, err
	}
	s3minio, err := NewS3(ctx, databases.Minio, cfg)
	if err != nil {
		l.Info("s3minio", "err", err.Error())
		return nil, err
	}
	httpsrv := httpserver.NewHTTPServer(cfg.HttpServer.Host, cfg.HttpServer.Port, cfg.HttpServer.Timeout,
		cfg.HttpServer.IdleTimeout,
	)
	authClient, err := authgrpc.NewAuthClient(l, cfg.Clients.Auth)
	if err != nil {
		l.Info("authClient", "err", err.Error())
		return nil, err
	}
	repo := NewRepository(databases, l)
	services := NewServices(repo, s3minio, authClient, l)
	handler := NewHandlers(l, services, httpsrv.Upgrader())

	return &App{
		Handlers:     handler,
		Cfg:          cfg,
		HTTPServer:   httpsrv,
		Log:          l,
		CancelLogger: cancelogger,
	}, nil
}
