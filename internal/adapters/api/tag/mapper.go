package tag

import (
	"github.com/google/uuid"
	"task-api/internal/domain/entities"
	"time"
)

func (t *CreateTagRequest) ToEntity() *entities.Tag {
	return &entities.Tag{
		Title:     t.Title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (t *UpdateTagRequest) ToEntity(ID uuid.UUID) *entities.Tag {
	return &entities.Tag{
		ID:        ID,
		Title:     t.Title,
		UpdatedAt: time.Now(),
	}
}

func FromEntityTag(e *entities.Tag) *TagResponse {
	return &TagResponse{
		ID:        e.ID,
		Title:     e.Title,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
