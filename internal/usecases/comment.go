package usecases

import (
	"context"
	"github.com/google/uuid"
	"task-api/internal/adapters/models"
	"task-api/internal/domain/entities"
	"task-api/internal/domain/repositories"
)

type CommentUseCase interface {
	GetAll(ctx context.Context) ([]*models.CommentWish, error)
	Create(ctx context.Context, comment *entities.Comment) (*models.CommentWish, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.CommentWish, error)
	Update(ctx context.Context, comment *entities.Comment) (*models.CommentWish, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type commentUseCase struct {
	repo repositories.CommentRepository
}

func NewCommentUseCase(repo repositories.CommentRepository) CommentUseCase {
	return &commentUseCase{repo: repo}
}

func (c *commentUseCase) GetAll(ctx context.Context) ([]*models.CommentWish, error) {
	return c.repo.GetAll(ctx)
}

func (c *commentUseCase) Create(ctx context.Context, comment *entities.Comment) (*models.CommentWish, error) {
	if err := c.repo.Create(ctx, comment); err != nil {
		return nil, err
	}
	res, err := c.repo.GetByID(ctx, comment.ID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *commentUseCase) GetByID(ctx context.Context, id uuid.UUID) (*models.CommentWish, error) {
	return c.repo.GetByID(ctx, id)
}

func (c *commentUseCase) Update(ctx context.Context, comment *entities.Comment) (*models.CommentWish, error) {
	if err := c.repo.Update(ctx, comment); err != nil {
		return nil, err
	}
	res, err := c.repo.GetByID(ctx, comment.ID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *commentUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return c.repo.Delete(ctx, id)
}
