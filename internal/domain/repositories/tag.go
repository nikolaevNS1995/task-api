package repositories

import (
	"context"
	"github.com/google/uuid"
	"task-api/internal/domain/entities"
)

type TagRepository interface {
	GetAllTags(ctx context.Context) ([]*entities.Tag, error)
	CreateTag(ctx context.Context, tag *entities.Tag) error
	GetTagByID(ctx context.Context, id uuid.UUID) (*entities.Tag, error)
	UpdateTag(ctx context.Context, tag *entities.Tag) error
	DeleteTag(ctx context.Context, id uuid.UUID) error
}
