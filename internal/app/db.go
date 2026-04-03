package app

import (
	"context"
	"github.com/Sanchir01/fitnow/internal/config"
	"github.com/Sanchir01/fitnow/pkg/db/connect"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type Database struct {
	PrimaryDB *pgxpool.Pool
}

func NewDataBases(cfg *config.Config, log *slog.Logger) (*Database, error) {
	pgxdb, err := connect.PGXNew(cfg, context.Background())
	if err != nil {
		log.Error("pgx connect error", err.Error())
		return nil, err
	}
	//redisdb, err := connect.RedisConnect(context.Background(), cfg.RedisDB.Host, cfg.RedisDB.Port,
	//	os.Getenv("REDIS_PASSWORD"), cfg.Env,
	//	cfg.RedisDB.DBNumber, cfg.RedisDB.Retries)
	//if err != nil {
	//	log.Error("redis connect error", err.Error())
	//	return nil, err
	//}
	return &Database{PrimaryDB: pgxdb}, nil
}

func (databases *Database) Close() error {
	databases.PrimaryDB.Close()

	return nil
}
