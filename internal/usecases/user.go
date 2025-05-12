package usecases

import (
	"context"
	"github.com/google/uuid"
	"task-api/internal/domain/entities"
	"task-api/internal/domain/repositories"
)

type UserUseCase interface {
	GetById(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userUseCase struct {
	repo repositories.UserRepository
}

func NewUserUseCase(repo repositories.UserRepository) UserUseCase {
	return &userUseCase{repo: repo}
}

func (u *userUseCase) GetById(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	return u.repo.GetById(ctx, id)
}

func (u *userUseCase) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	return u.repo.GetByEmail(ctx, email)
}

func (u *userUseCase) Update(ctx context.Context, user *entities.User) error {
	if err := u.repo.Update(ctx, user); err != nil {
		return err
	}
	user, err := u.repo.GetById(ctx, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
