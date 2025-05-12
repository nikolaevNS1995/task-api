package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"task-api/internal/domain/entities"
	"task-api/internal/domain/repositories"
)

type RefreshTokenPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewRefreshTokenPostgresRepository(pool *pgxpool.Pool) *RefreshTokenPostgresRepository {
	return &RefreshTokenPostgresRepository{pool: pool}
}

var _ repositories.RefreshTokenRepository = new(RefreshTokenPostgresRepository)

func (r *RefreshTokenPostgresRepository) Create(ctx context.Context, token *entities.RefreshToken) error {
	sql := `INSERT INTO users.refresh_tokens (token, user_id, expires_at, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, sql, token.Token, token.UserID, token.ExpiresAt, token.CreatedAt)
	return err
}

func (r *RefreshTokenPostgresRepository) GetByToken(ctx context.Context, tokenID string) (*entities.RefreshToken, error) {
	sql := `SELECT token, user_id, expires_at, created_at FROM users.refresh_tokens WHERE token = $1`
	row := r.pool.QueryRow(ctx, sql, tokenID)
	token := &entities.RefreshToken{}
	if err := row.Scan(&token.Token, &token.UserID, &token.ExpiresAt, &token.CreatedAt); err != nil {
		return nil, err
	}
	return token, nil
}

func (r *RefreshTokenPostgresRepository) Delete(ctx context.Context, tokenID string) error {
	sql := `DELETE FROM users.refresh_tokens WHERE token = $1`
	_, err := r.pool.Exec(ctx, sql, tokenID)
	return err
}
