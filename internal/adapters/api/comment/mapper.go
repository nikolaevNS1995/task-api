package comment

import (
	"github.com/google/uuid"
	"task-api/internal/adapters/models"
	"task-api/internal/domain/entities"
	"time"
)

func (r *CreateCommentRequest) ToEntity() *entities.Comment {
	return &entities.Comment{
		TaskID:    r.TaskID,
		Author:    r.AuthorID,
		Content:   r.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (r *UpdateCommentRequest) ToEntity(ID uuid.UUID) *entities.Comment {
	return &entities.Comment{
		ID:        ID,
		Content:   r.Content,
		UpdatedAt: time.Now(),
	}
}

func FromModelComment(m *models.CommentWish) *CommentResponse {
	return &CommentResponse{
		ID:     m.Comment.ID,
		TaskID: m.Comment.TaskID,
		Author: Author{
			ID:    m.Author.ID,
			Name:  m.Author.Name,
			Email: m.Author.Email,
		},
		Content:   m.Comment.Content,
		CreatedAt: m.Comment.CreatedAt,
		UpdatedAt: m.Comment.UpdatedAt,
	}
}
