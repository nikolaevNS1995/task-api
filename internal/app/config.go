package app

import (
	"log"
	"task-api/pkg/config"
)

func NewConfig() (*config.AppConfig, error) {
	var cfg config.AppConfig
	if err := cfg.ReadEnvConfig(); err != nil {
		log.Fatalf("reading config failed: %v", err)
		return nil, err
	}
	if err := cfg.Validate(); err != nil {
		log.Fatalf("validate config failed: %v", err)
		return nil, err
	}
	return &cfg, nil
}
