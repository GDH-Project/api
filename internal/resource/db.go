package resource

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func InitDB(connectionString string, log *zap.Logger) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		log.Fatal("failed to parse DB config", zap.Error(err))
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("failed to connect to DB", zap.Error(err))
	}

	if err = pool.Ping(context.Background()); err != nil {
		log.Fatal("failed to ping DB", zap.Error(err))
	}

	log.Info("Database connection pool initialized successfully")

	return pool
}
