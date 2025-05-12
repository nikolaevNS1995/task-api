package repositories

import (
	"context"
	"github.com/google/uuid"
	"task-api/internal/adapters/models"
	"task-api/internal/domain/entities"
)

type CommentRepository interface {
	GetAll(ctx context.Context) ([]*models.CommentWish, error)
	Create(ctx context.Context, tag *entities.Comment) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.CommentWish, error)
	Update(ctx context.Context, comment *entities.Comment) error
	Delete(ctx context.Context, id uuid.UUID) error
}
