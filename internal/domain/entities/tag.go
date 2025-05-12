package entities

import (
	"github.com/google/uuid"
	"time"
)

type Tag struct {
	ID        uuid.UUID
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
