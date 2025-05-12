package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"task-api/internal/adapters/models"
	"task-api/internal/domain/entities"
	"task-api/internal/domain/repositories"
)

type CommentRepository struct {
	pool *pgxpool.Pool
}

var _ repositories.CommentRepository = new(CommentRepository)

func NewCommentRepository(pool *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{pool: pool}
}

func (c *CommentRepository) GetAll(ctx context.Context) ([]*models.CommentWish, error) {
	sql := `SELECT c.id, c.task_id, c.author_id, c.content, c.created_at, c.updated_at, u.id, u.name, u.email 
			FROM tasks.comments c
			JOIN users.users u ON u.id = c.author_id`
	rows, err := c.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []*models.CommentWish
	for rows.Next() {
		row := &models.CommentWish{}
		if err := rows.Scan(
			&row.Comment.ID,
			&row.Comment.TaskID,
			&row.Comment.Author,
			&row.Comment.Content,
			&row.Comment.CreatedAt,
			&row.Comment.UpdatedAt,
			&row.Author.ID,
			&row.Author.Name,
			&row.Author.Email,
		); err != nil {
			return nil, err
		}
		comments = append(comments, row)
	}
	return comments, nil
}

func (c *CommentRepository) Create(ctx context.Context, comment *entities.Comment) error {
	sql := `INSERT INTO tasks.comments (task_id, author_id, content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return c.pool.QueryRow(ctx, sql, comment.TaskID, comment.Author, comment.Content, comment.CreatedAt, comment.UpdatedAt).Scan(&comment.ID)
}

func (c *CommentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.CommentWish, error) {
	sql := `SELECT c.id, c.task_id, c.author_id, c.content, c.created_at, c.updated_at, u.id, u.name, u.email
			FROM tasks.comments c
			JOIN users.users u ON u.id = c.author_id
			WHERE c.id = $1`
	row := c.pool.QueryRow(ctx, sql, id)
	res := &models.CommentWish{}
	if err := row.Scan(
		&res.Comment.ID,
		&res.Comment.TaskID,
		&res.Comment.Author,
		&res.Comment.Content,
		&res.Comment.CreatedAt,
		&res.Comment.UpdatedAt,
		&res.Author.ID,
		&res.Author.Name,
		&res.Author.Email,
	); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *CommentRepository) Update(ctx context.Context, comment *entities.Comment) error {
	sql := `UPDATE tasks.comments SET content = $1, updated_at = $2 WHERE id = $3`
	_, err := c.pool.Exec(ctx, sql, comment.Content, comment.UpdatedAt, comment.ID)
	return err
}

func (c *CommentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `DELETE FROM tasks.comments WHERE id = $1`
	_, err := c.pool.Exec(ctx, sql, id)
	return err
}
