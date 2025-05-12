package task

import "github.com/google/uuid"

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type TagRequest struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Status      string `json:"status" binding:"required"`
}
