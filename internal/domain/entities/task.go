package entities

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	ID          uuid.UUID
	Title       string
	Description string
	Status      string
	CreatedBy   uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
