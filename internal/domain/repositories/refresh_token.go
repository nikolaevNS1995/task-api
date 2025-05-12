package repositories

import (
	"context"
	"task-api/internal/domain/entities"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *entities.RefreshToken) error
	GetByToken(ctx context.Context, tokenID string) (*entities.RefreshToken, error)
	Delete(ctx context.Context, tokenID string) error
}
