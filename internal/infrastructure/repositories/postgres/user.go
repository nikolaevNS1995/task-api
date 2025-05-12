package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"task-api/internal/domain/entities"
	"task-api/internal/domain/repositories"
	"time"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

var _ repositories.UserRepository = new(UserRepository)

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (u *UserRepository) Create(ctx context.Context, user *entities.User) error {
	sql := `INSERT INTO users.users (name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return u.pool.QueryRow(ctx, sql, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
}

func (u *UserRepository) GetById(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	sql := `SELECT id, name, email, password, created_at, updated_at FROM users.users WHERE id = $1`
	row := u.pool.QueryRow(ctx, sql, id)
	user := &entities.User{}
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	sql := `SELECT id, name, email, password, created_at, updated_at FROM users.users WHERE email = $1`
	row := u.pool.QueryRow(ctx, sql, email)
	user := &entities.User{}
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) Update(ctx context.Context, user *entities.User) error {
	sql := `UPDATE users.users 
			SET name = $1, email = $2, password = $3, updated_at = $4
			WHERE id = $5`
	user.UpdatedAt = time.Now()
	_, err := u.pool.Exec(ctx, sql, user.Name, user.Email, user.Password, user.UpdatedAt, user.ID)
	return err
}

func (u *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `DELETE FROM users.users WHERE id = $1`
	_, err := u.pool.Exec(ctx, sql, id)
	return err
}
