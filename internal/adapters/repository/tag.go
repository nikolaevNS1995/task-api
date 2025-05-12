package repository

import (
	"github.com/google/uuid"
	"task-api/internal/domain/entities"
	"time"
)

type Tag struct {
	ID        uuid.UUID `db:"id"`
	Title     string    `db:"title"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (t *Tag) ToEntity() *entities.Tag {
	return &entities.Tag{
		ID:        t.ID,
		Title:     t.Title,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func FromEntityTag(e *entities.Tag) *Tag {
	return &Tag{
		ID:        e.ID,
		Title:     e.Title,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
