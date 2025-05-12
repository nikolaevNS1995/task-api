package usecases

import (
	"context"
	"github.com/google/uuid"
	"task-api/internal/domain/entities"
	"task-api/internal/domain/repositories"
)

type TagUseCase interface {
	Create(ctx context.Context, tag *entities.Tag) (*entities.Tag, error)
	GetTag(ctx context.Context, id uuid.UUID) (*entities.Tag, error)
	GetTags(ctx context.Context) ([]*entities.Tag, error)
	Update(ctx context.Context, tag *entities.Tag) error
	Delete(ctx context.Context, id uuid.UUID) error
}
type tagUseCase struct {
	repo repositories.TagRepository
}

func NewTagsUseCase(repo repositories.TagRepository) TagUseCase {
	return &tagUseCase{repo: repo}
}

func (t *tagUseCase) Create(ctx context.Context, tag *entities.Tag) (*entities.Tag, error) {
	if err := t.repo.CreateTag(ctx, tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (t *tagUseCase) GetTag(ctx context.Context, id uuid.UUID) (*entities.Tag, error) {
	return t.repo.GetTagByID(ctx, id)
}

func (t *tagUseCase) GetTags(ctx context.Context) ([]*entities.Tag, error) {
	return t.repo.GetAllTags(ctx)
}

func (t *tagUseCase) Update(ctx context.Context, tag *entities.Tag) error {
	if err := t.repo.UpdateTag(ctx, tag); err != nil {
		return err
	}
	tag, err := t.repo.GetTagByID(ctx, tag.ID)
	if err != nil {
		return err
	}
	return nil
}

func (t *tagUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return t.repo.DeleteTag(ctx, id)
}
