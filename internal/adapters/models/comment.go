package models

import "task-api/internal/domain/entities"

type CommentWish struct {
	Comment entities.Comment
	Author  entities.User
}
