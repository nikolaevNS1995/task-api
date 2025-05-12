package repositories

import (
	"context"
	"github.com/google/uuid"
	"task-api/internal/adapters/models"
	"task-api/internal/domain/entities"
)

type TaskRepository interface {
	GetAllTasks(ctx context.Context) ([]*models.Task, error)
	GetAllTasksByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Task, error)
	CreateTask(ctx context.Context, task *entities.Task) error
	GetTaskByID(ctx context.Context, id uuid.UUID) (*models.Task, error)
	UpdateTask(ctx context.Context, task *entities.Task) error
	DeleteTask(ctx context.Context, id uuid.UUID) error

	AddTags(ctx context.Context, taskID, tagID uuid.UUID) error
	RemoveTags(ctx context.Context, taskID, tagID uuid.UUID) error
	GetTags(ctx context.Context, taskID uuid.UUID) ([]*entities.Tag, error)
	GetTagsForManyTasks(ctx context.Context, taskIDs []uuid.UUID) ([]*models.TagWishTaskID, error)
	GetComments(ctx context.Context, taskID uuid.UUID) ([]*models.CommentWish, error)
}
