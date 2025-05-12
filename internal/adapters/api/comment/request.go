package comment

import "github.com/google/uuid"

type CreateCommentRequest struct {
	TaskID   uuid.UUID `json:"task_id" binding:"required"`
	AuthorID uuid.UUID `json:"author" binding:"required"`
	Content  string    `json:"content" binding:"required"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}
