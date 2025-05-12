package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"task-api/internal/adapters/models"
	"task-api/internal/domain/entities"
	"task-api/internal/domain/repositories"
)

type TaskRepository struct {
	pool *pgxpool.Pool
}

var _ repositories.TaskRepository = new(TaskRepository)

func NewTaskPostgresRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{pool: pool}
}

func (r *TaskRepository) GetAllTasks(ctx context.Context) ([]*models.Task, error) {
	sql := `SELECT t.id, t.title, t.description, t.status, t.created_by, t.created_at, t.updated_at, u.id, u.name, u.email
			FROM tasks.tasks t
			JOIN users.users u ON u.id = t.created_by`
	rows, err := r.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		if err := rows.Scan(
			&task.Task.ID,
			&task.Task.Title,
			&task.Task.Description,
			&task.Task.Status,
			&task.Task.CreatedBy,
			&task.Task.CreatedAt,
			&task.Task.UpdatedAt,
			&task.User.ID,
			&task.User.Name,
			&task.User.Email,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepository) GetAllTasksByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Task, error) {
	sql := `SELECT id, title, description, status, created_by, created_at, updated_at
			FROM tasks.tasks WHERE created_by = $1`
	rows, err := r.pool.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*entities.Task
	for rows.Next() {
		task := &entities.Task{}
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedBy,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepository) CreateTask(ctx context.Context, task *entities.Task) error {
	sql := `INSERT INTO tasks.tasks (title, description, status, created_by, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	return r.pool.QueryRow(ctx, sql, task.Title, task.Description, task.Status, task.CreatedBy, task.CreatedAt, task.UpdatedAt).Scan(&task.ID)
}

func (r *TaskRepository) GetTaskByID(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	sql := `SELECT t.id, t.title, t.description, t.status, t.created_by, t.created_at, t.updated_at, u.id, u.name, u.email
			FROM tasks.tasks t
			JOIN users.users u ON u.id = t.created_by
			WHERE t.id = $1`
	row := r.pool.QueryRow(ctx, sql, id)
	task := &models.Task{}
	if err := row.Scan(
		&task.Task.ID,
		&task.Task.Title,
		&task.Task.Description,
		&task.Task.Status,
		&task.Task.CreatedBy,
		&task.Task.CreatedAt,
		&task.Task.UpdatedAt,
		&task.User.ID,
		&task.User.Name,
		&task.User.Email,
	); err != nil {
		return nil, err
	}
	return task, nil

}

func (r *TaskRepository) UpdateTask(ctx context.Context, task *entities.Task) error {
	sql := `UPDATE tasks.tasks 
			SET title = $1, description = $2, status = $3, updated_at = $4 
			WHERE id = $5`
	_, err := r.pool.Exec(ctx, sql, task.Title, task.Description, task.Status, task.UpdatedAt, task.ID)
	return err
}

func (r *TaskRepository) DeleteTask(ctx context.Context, id uuid.UUID) error {
	sql := `DELETE FROM tasks.tasks WHERE id = $1`
	_, err := r.pool.Exec(ctx, sql, id)
	return err
}

func (r *TaskRepository) AddTags(ctx context.Context, taskID, tagID uuid.UUID) error {
	sql := `INSERT INTO tasks.tasks_tags (task_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.pool.Exec(ctx, sql, taskID, tagID)
	return err
}

func (r *TaskRepository) RemoveTags(ctx context.Context, taskID, tagID uuid.UUID) error {
	sql := `DELETE FROM tasks.tasks_tags WHERE task_id = $1 AND tag_id = $2`
	_, err := r.pool.Exec(ctx, sql, taskID, tagID)
	return err
}

func (r *TaskRepository) GetTags(ctx context.Context, taskID uuid.UUID) ([]*entities.Tag, error) {
	sql := `SELECT t.id, t.title 
			FROM tasks.tags t
			JOIN tasks.tasks_tags tt ON t.id = tt.tag_id
			WHERE tt.task_id = $1`
	rows, err := r.pool.Query(ctx, sql, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tags []*entities.Tag
	for rows.Next() {
		tag := &entities.Tag{}
		if err := rows.Scan(&tag.ID, &tag.Title); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (r *TaskRepository) GetTagsForManyTasks(ctx context.Context, taskIDs []uuid.UUID) ([]*models.TagWishTaskID, error) {
	sql := `SELECT tt.task_id, t.id, t.title 
			FROM tasks.tags t
			LEFT JOIN tasks.tasks_tags tt ON t.id = tt.tag_id
			WHERE tt.task_id = ANY($1) OR tt.task_id IS NULL`
	rows, err := r.pool.Query(ctx, sql, taskIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tags []*models.TagWishTaskID
	for rows.Next() {
		tag := &models.TagWishTaskID{}
		if err := rows.Scan(&tag.TaskID, &tag.Tag.ID, &tag.Tag.Title); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (r *TaskRepository) GetComments(ctx context.Context, taskID uuid.UUID) ([]*models.CommentWish, error) {
	sql := `SELECT c.id, c.task_id, c.author_id, c.content, c.created_at, c.updated_at, u.id, u.name, u.email 
			FROM tasks.comments c
			JOIN users.users u ON u.id = c.author_id
			WHERE c.task_id = $1`
	rows, err := r.pool.Query(ctx, sql, taskID)
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
