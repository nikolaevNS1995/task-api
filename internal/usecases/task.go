package usecases

import (
	"context"
	"github.com/google/uuid"
	"task-api/internal/adapters/models"
	"task-api/internal/domain/entities"
	"task-api/internal/domain/repositories"
)

type TaskUseCase interface {
	Create(ctx context.Context, task *entities.Task) (*models.Task, error)
	GetTask(ctx context.Context, id uuid.UUID) (*models.Task, error)
	GetTasks(ctx context.Context) ([]*models.Task, error)
	GetTasksByUserID(ctx context.Context, userID uuid.UUID) ([]*models.TasksWishTags, error)
	Update(ctx context.Context, task *entities.Task) (*models.Task, error)
	Delete(ctx context.Context, id uuid.UUID) error
	AddTags(ctx context.Context, taskID uuid.UUID, tags []*entities.Tag) error
	RemoveTags(ctx context.Context, taskID uuid.UUID, tags []*entities.Tag) error
}

type tasksUseCase struct {
	repo repositories.TaskRepository
}

func NewTasksUseCase(repo repositories.TaskRepository) TaskUseCase {
	return &tasksUseCase{repo: repo}
}

func (t *tasksUseCase) Create(ctx context.Context, task *entities.Task) (*models.Task, error) {
	if err := t.repo.CreateTask(ctx, task); err != nil {
		return nil, err
	}
	model, err := t.repo.GetTaskByID(ctx, task.ID)
	if err != nil {
		return nil, err
	}

	comments, err := t.repo.GetComments(ctx, task.ID)
	if err != nil {
		return nil, err
	}
	for _, comment := range comments {
		model.Comments = append(model.Comments, *comment)
	}

	tags, err := t.repo.GetTags(ctx, task.ID)
	if err != nil {
		return nil, err
	}
	for _, tag := range tags {
		model.Tags = append(model.Tags, *tag)
	}
	return model, nil
}

func (t *tasksUseCase) GetTask(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	task, err := t.repo.GetTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}

	comments, err := t.repo.GetComments(ctx, task.Task.ID)
	if err != nil {
		return nil, err
	}
	for _, comment := range comments {
		task.Comments = append(task.Comments, *comment)
	}

	tags, err := t.repo.GetTags(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, tag := range tags {
		task.Tags = append(task.Tags, *tag)
	}
	return task, nil
}

func (t *tasksUseCase) GetTasksByUserID(ctx context.Context, userID uuid.UUID) ([]*models.TasksWishTags, error) {
	tasks, err := t.repo.GetAllTasksByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	var taskIds []uuid.UUID
	for _, task := range tasks {
		taskIds = append(taskIds, task.ID)
	}
	tags, err := t.repo.GetTagsForManyTasks(ctx, taskIds)
	if err != nil {
		return nil, err
	}
	tagsMap := make(map[uuid.UUID][]entities.Tag)
	for _, tag := range tags {
		tagsMap[tag.TaskID] = append(tagsMap[tag.TaskID], entities.Tag{
			ID:    tag.Tag.ID,
			Title: tag.Tag.Title,
		})
	}
	var tasksWishTags []*models.TasksWishTags
	for _, task := range tasks {
		taskWishTags := &models.TasksWishTags{
			Task: *task,
			Tags: tagsMap[task.ID],
		}
		tasksWishTags = append(tasksWishTags, taskWishTags)
	}
	return tasksWishTags, nil
}

func (t *tasksUseCase) GetTasks(ctx context.Context) ([]*models.Task, error) {
	tasks, err := t.repo.GetAllTasks(ctx)
	if err != nil {
		return nil, err
	}
	var taskIds []uuid.UUID
	for _, task := range tasks {
		taskIds = append(taskIds, task.Task.ID)
	}
	tags, err := t.repo.GetTagsForManyTasks(ctx, taskIds)
	if err != nil {
		return nil, err
	}
	tagsMap := make(map[uuid.UUID][]entities.Tag)
	for _, tag := range tags {
		tagsMap[tag.TaskID] = append(tagsMap[tag.TaskID], entities.Tag{
			ID:    tag.Tag.ID,
			Title: tag.Tag.Title,
		})
	}

	for _, task := range tasks {
		task.Tags = tagsMap[task.Task.ID]
	}
	return tasks, nil
}

func (t *tasksUseCase) Update(ctx context.Context, task *entities.Task) (*models.Task, error) {
	if err := t.repo.UpdateTask(ctx, task); err != nil {
		return nil, err
	}
	model, err := t.repo.GetTaskByID(ctx, task.ID)
	if err != nil {
		return nil, err
	}

	comments, err := t.repo.GetComments(ctx, task.ID)
	if err != nil {
		return nil, err
	}
	for _, comment := range comments {
		model.Comments = append(model.Comments, *comment)
	}

	tags, err := t.repo.GetTags(ctx, task.ID)
	if err != nil {
		return nil, err
	}
	for _, tag := range tags {
		model.Tags = append(model.Tags, *tag)
	}
	return model, err
}

func (t *tasksUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return t.repo.DeleteTask(ctx, id)
}

func (t *tasksUseCase) AddTags(ctx context.Context, taskID uuid.UUID, tags []*entities.Tag) error {
	for _, tag := range tags {
		if err := t.repo.AddTags(ctx, taskID, tag.ID); err != nil {
			return err
		}
	}
	return nil
}

func (t *tasksUseCase) RemoveTags(ctx context.Context, taskID uuid.UUID, tags []*entities.Tag) error {
	for _, tag := range tags {
		if err := t.repo.RemoveTags(ctx, taskID, tag.ID); err != nil {
			return err
		}
	}
	return nil
}
