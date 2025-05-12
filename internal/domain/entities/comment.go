package entities

import (
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	ID        uuid.UUID
	TaskID    uuid.UUID
	Author    uuid.UUID
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
