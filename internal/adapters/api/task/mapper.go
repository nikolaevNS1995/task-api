package task

import (
	"github.com/google/uuid"
	commentRes "task-api/internal/adapters/api/comment"
	"task-api/internal/adapters/models"
	"task-api/internal/domain/entities"
	"time"
)

func (req *CreateTaskRequest) ToEntity(userID uuid.UUID) *entities.Task {
	return &entities.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      "new",
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (req *UpdateTaskRequest) ToEntity(ID uuid.UUID) *entities.Task {
	return &entities.Task{
		ID:          ID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		UpdatedAt:   time.Now(),
	}
}

func FromModelTask(m *models.Task) *TaskResponse {
	res := &TaskResponse{
		ID:          m.Task.ID,
		Title:       m.Task.Title,
		Description: m.Task.Description,
		Status:      m.Task.Status,
		CreatedBy: Creator{
			ID:    m.User.ID,
			Name:  m.User.Name,
			Email: m.User.Email,
		},
		CreatedAt: m.Task.CreatedAt,
		UpdatedAt: m.Task.UpdatedAt,
	}
	for _, comment := range m.Comments {
		res.Comments = append(res.Comments,
			commentRes.CommentResponse{
				ID:     comment.Comment.ID,
				TaskID: comment.Comment.TaskID,
				Author: commentRes.Author{
					ID:    comment.Author.ID,
					Name:  comment.Author.Name,
					Email: comment.Author.Email,
				},
				Content:   comment.Comment.Content,
				CreatedAt: comment.Comment.CreatedAt,
				UpdatedAt: comment.Comment.UpdatedAt,
			})
	}
	for _, tag := range m.Tags {
		res.Tags = append(res.Tags, Tags{ID: tag.ID, Title: tag.Title})
	}
	return res
}

func FromModelTaskForAll(m *models.TasksWishTags) *TaskAllResponse {
	res := &TaskAllResponse{
		ID:          m.Task.ID,
		Title:       m.Task.Title,
		Description: m.Task.Description,
		Status:      m.Task.Status,
		CreatedAt:   m.Task.CreatedAt,
		UpdatedAt:   m.Task.UpdatedAt,
	}
	for _, tag := range m.Tags {
		res.Tags = append(res.Tags, Tags{ID: tag.ID, Title: tag.Title})
	}
	return res
}

func (r *TagRequest) ToEntity() *entities.Tag {
	return &entities.Tag{
		ID: r.ID,
	}
}
