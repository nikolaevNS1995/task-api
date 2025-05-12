package user

import (
	"github.com/google/uuid"
	"task-api/internal/domain/entities"
	"time"
)

func (u *CreateUserRequest) ToEntity() *entities.User {
	return &entities.User{
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (u *UpdateUserRequest) ToEntity(ID uuid.UUID) *entities.User {
	return &entities.User{
		ID:        ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		UpdatedAt: time.Now(),
	}
}

func FromEntityUser(e *entities.User) *UserResponse {
	return &UserResponse{
		ID:        e.ID,
		Name:      e.Name,
		Email:     e.Email,
		Password:  e.Password,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
