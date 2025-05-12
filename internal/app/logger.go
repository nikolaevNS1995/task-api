package app

import (
	"go.uber.org/zap"
	"log"
	"task-api/pkg/config"
	"task-api/pkg/logger"
)

func NewLogger(cfg *config.AppConfig) (*zap.Logger, error) {
	l, err := logger.CreateLogger(cfg.Logger)
	if err != nil {
		log.Fatalf("logger init failed: %v", err)
		return nil, err
	}
	zap.ReplaceGlobals(l)
	return l, nil
}
