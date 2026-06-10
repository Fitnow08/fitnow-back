package app

import (
	"context"
	"fmt"
	"github.com/Sanchir01/fitnow/internal/config"
	"github.com/Sanchir01/fitnow/internal/feature/events"
	constants "github.com/Sanchir01/fitnow/internal/models/contants"
	httpserver "github.com/Sanchir01/fitnow/internal/servers/http"
	"github.com/Sanchir01/fitnow/pkg/kafka"
	"github.com/Sanchir01/fitnow/pkg/logger"
	kafkago "github.com/segmentio/kafka-go"
	"log/slog"
)

type App struct {
	Handlers      *Handlers
	Cfg           *config.Config
	HTTPServer    *httpserver.Server
	Log           *slog.Logger
	KafkaProducer *kafka.Producer
	EventService  *events.Service
	CancelLogger  func()
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
	clients, err := NewClients(cfg, l)
	if err != nil {
		l.Info("clients", "err", err.Error())
		return nil, err
	}
	producer := kafka.NewProducer(cfg.Kafka.Brokers, l)
	if err := producer.EnsureTopics(ctx, kafkago.TopicConfig{
		Topic:             constants.TopicProgramCreated,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}); err != nil {
		l.Error("kafka ensure topics", slog.String("err", err.Error()))
		return nil, err
	}
	repo := NewRepository(databases, l)
	services := NewServices(repo, s3minio, clients, producer, l)
	handler := NewHandlers(l, services, httpsrv.Upgrader())

	return &App{
		Handlers:      handler,
		Cfg:           cfg,
		HTTPServer:    httpsrv,
		Log:           l,
		KafkaProducer: producer,
		CancelLogger:  cancelogger,
		EventService:  services.EventService,
	}, nil
}
