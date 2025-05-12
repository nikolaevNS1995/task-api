package app

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"task-api/pkg/config"
	"task-api/pkg/connectors"
)

func NewPostgresConnection(lc fx.Lifecycle, cfg *config.AppConfig, logger *zap.Logger) (*connectors.PostgresConnect, error) {
	pool, err := connectors.NewPostgresConnect(&cfg.MainStorage.Postgres)
	if err != nil {
		logger.Error("failed to connect to PostgreSQL", zap.Error(err))
		return nil, err
	}
	logger.Info("connected to PostgreSQL")

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("closing PostgreSQL connection")
			if err := pool.Close(); err != nil {
				logger.Error("failed to close PostgreSQL connection", zap.Error(err))
			}
			logger.Info("PostgreSQL connection closed")
			return nil
		},
	})
	return pool, nil
}
