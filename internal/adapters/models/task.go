package models

import (
	"github.com/google/uuid"
	"task-api/internal/domain/entities"
)

type Task struct {
	Task     entities.Task
	Tags     []entities.Tag
	User     entities.User
	Comments []CommentWish
}

type TasksWishTags struct {
	Task entities.Task
	Tags []entities.Tag
}

type TagWishTaskID struct {
	TaskID uuid.UUID
	Tag    entities.Tag
}
