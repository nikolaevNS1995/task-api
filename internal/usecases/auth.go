package usecases

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"task-api/internal/domain/entities"
	"task-api/internal/domain/repositories"
	"time"
)

type AuthUseCase interface {
	Login(ctx context.Context, email, password string) (*entities.User, error)
	Register(ctx context.Context, user *entities.User) (*entities.User, error)
	CreateRefreshToken(ctx context.Context, userID uuid.UUID, ExpiresAt time.Time) (*entities.RefreshToken, error)
	GetRefreshToken(ctx context.Context, tokenID string) (*entities.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
}

type authUseCase struct {
	repoUser   repositories.UserRepository
	repoRefTok repositories.RefreshTokenRepository
}

func NewAuthUseCase(repoUser repositories.UserRepository, repoRefTok repositories.RefreshTokenRepository) AuthUseCase {
	return &authUseCase{repoUser: repoUser, repoRefTok: repoRefTok}
}

func (a *authUseCase) Login(ctx context.Context, email, password string) (*entities.User, error) {
	user, err := a.repoUser.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (a *authUseCase) Register(ctx context.Context, user *entities.User) (*entities.User, error) {
	existing, _ := a.repoUser.GetByEmail(ctx, user.Email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	if err := a.repoUser.Create(ctx, user); err != nil {
		return nil, err
	}
	create, err := a.repoUser.GetById(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return create, nil
}

func (a *authUseCase) CreateRefreshToken(ctx context.Context, userID uuid.UUID, ExpiresAt time.Time) (*entities.RefreshToken, error) {
	refreshToken := &entities.RefreshToken{
		Token:     uuid.New(),
		UserID:    userID,
		ExpiresAt: ExpiresAt,
		CreatedAt: time.Now(),
	}
	if err := a.repoRefTok.Create(ctx, refreshToken); err != nil {
		return nil, err
	}
	return refreshToken, nil

}

func (a *authUseCase) GetRefreshToken(ctx context.Context, tokenID string) (*entities.RefreshToken, error) {
	return a.repoRefTok.GetByToken(ctx, tokenID)
}

func (a *authUseCase) DeleteRefreshToken(ctx context.Context, tokenID string) error {
	return a.repoRefTok.Delete(ctx, tokenID)
}
