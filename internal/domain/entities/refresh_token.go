package entities

import (
	"github.com/google/uuid"
	"time"
)

type RefreshToken struct {
	Token     uuid.UUID
	UserID    uuid.UUID
	ExpiresAt time.Time
	CreatedAt time.Time
}
