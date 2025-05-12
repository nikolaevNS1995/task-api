package app

import (
	"task-api/internal/infrastructure/repositories/postgres"
	"task-api/pkg/connectors"
)

type Repositories struct {
	taskRepo         *postgres.TaskRepository
	tagRepo          *postgres.TagRepository
	commentRepo      *postgres.CommentRepository
	userRepo         *postgres.UserRepository
	refreshTokenRepo *postgres.RefreshTokenPostgresRepository
}

func NewRopositories(pool *connectors.PostgresConnect) *Repositories {
	return &Repositories{
		taskRepo:         postgres.NewTaskPostgresRepository(pool.Pool),
		tagRepo:          postgres.NewTagPostgresRepository(pool.Pool),
		commentRepo:      postgres.NewCommentRepository(pool.Pool),
		userRepo:         postgres.NewUserRepository(pool.Pool),
		refreshTokenRepo: postgres.NewRefreshTokenPostgresRepository(pool.Pool),
	}
}
