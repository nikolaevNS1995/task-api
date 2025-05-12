package task

import (
	"github.com/google/uuid"
	"task-api/internal/adapters/api/comment"
	"time"
)

type TaskResponse struct {
	ID          uuid.UUID                 `json:"id"`
	Title       string                    `json:"title"`
	Description string                    `json:"description"`
	Status      string                    `json:"status"`
	CreatedBy   Creator                   `json:"created_by"`
	Tags        []Tags                    `json:"tags"`
	Comments    []comment.CommentResponse `json:"comments"`
	CreatedAt   time.Time                 `json:"created_at"`
	UpdatedAt   time.Time                 `json:"updated_at"`
}

type TaskAllResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Tags        []Tags    `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Creator struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type Tags struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}
