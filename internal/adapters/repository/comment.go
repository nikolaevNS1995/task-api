package repository

import (
	"github.com/google/uuid"
	"task-api/internal/domain/entities"
	"time"
)

type Comment struct {
	ID        uuid.UUID `db:"id"`
	TaskID    uuid.UUID `db:"task_id"`
	Author    uuid.UUID `db:"author"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (c *Comment) ToEntity() *entities.Comment {
	return &entities.Comment{
		ID:        c.ID,
		TaskID:    c.TaskID,
		Author:    c.Author,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func FromEntityComment(e *entities.Comment) *Comment {
	return &Comment{
		ID:        e.ID,
		TaskID:    e.TaskID,
		Author:    e.Author,
		Content:   e.Content,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
