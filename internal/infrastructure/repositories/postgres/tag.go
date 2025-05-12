package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"task-api/internal/domain/entities"
	"task-api/internal/domain/repositories"
	"time"
)

type TagRepository struct {
	pool *pgxpool.Pool
}

var _ repositories.TagRepository = new(TagRepository)

func NewTagPostgresRepository(pool *pgxpool.Pool) *TagRepository {
	return &TagRepository{pool: pool}
}

func (t *TagRepository) GetAllTags(ctx context.Context) ([]*entities.Tag, error) {
	sql := `SELECT id, title, created_at, updated_at FROM tasks.tags`
	rows, err := t.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*entities.Tag
	for rows.Next() {
		tag := &entities.Tag{}
		if err := rows.Scan(&tag.ID, &tag.Title, &tag.CreatedAt, &tag.UpdatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (t *TagRepository) CreateTag(ctx context.Context, tag *entities.Tag) error {
	sql := `INSERT INTO tasks.tags (title, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id`
	return t.pool.QueryRow(ctx, sql, tag.Title, tag.CreatedAt, tag.UpdatedAt).Scan(&tag.ID)
}

func (t *TagRepository) GetTagByID(ctx context.Context, id uuid.UUID) (*entities.Tag, error) {
	sql := `SELECT id, title, created_at, updated_at FROM tasks.tags WHERE id = $1`
	row := t.pool.QueryRow(ctx, sql, id)
	tag := &entities.Tag{}
	if err := row.Scan(&tag.ID, &tag.Title, &tag.CreatedAt, &tag.UpdatedAt); err != nil {
		return nil, err
	}
	return tag, nil
}

func (t *TagRepository) UpdateTag(ctx context.Context, tag *entities.Tag) error {
	sql := `UPDATE tasks.tags SET title = $1, updated_at = $2 WHERE id = $3`
	tag.UpdatedAt = time.Now()
	_, err := t.pool.Exec(ctx, sql, tag.Title, tag.UpdatedAt, tag.ID)
	return err
}

func (t *TagRepository) DeleteTag(ctx context.Context, id uuid.UUID) error {
	sql := `DELETE FROM tasks.tags WHERE id = $1`
	_, err := t.pool.Exec(ctx, sql, id)
	return err
}
