package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"bench-hub/internal/model"
	"bench-hub/internal/repository"
)

type TaskRepo struct {
	pool *pgxpool.Pool
}

func NewTaskRepo(pool *pgxpool.Pool) *TaskRepo {
	return &TaskRepo{pool: pool}
}

func (r *TaskRepo) Create(ctx context.Context, task *model.Task) error {
	if task.ID == "" {
		task.ID = uuid.NewString()
	}

	row := r.pool.QueryRow(ctx,
		"INSERT INTO locust_tasks (id, name, script_id, users_count, spawn_rate, duration_seconds, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING created_at, updated_at",
		task.ID,
		task.Name,
		task.ScriptID,
		task.UsersCount,
		task.SpawnRate,
		task.DurationSeconds,
		task.Status,
	)

	return row.Scan(&task.CreatedAt, &task.UpdatedAt)
}

func (r *TaskRepo) GetByID(ctx context.Context, id string) (*model.Task, error) {
	task := &model.Task{}
	row := r.pool.QueryRow(ctx,
		"SELECT id, name, script_id, users_count, spawn_rate, duration_seconds, status, created_at, updated_at, started_at, finished_at FROM locust_tasks WHERE id = $1",
		id,
	)
	if err := row.Scan(
		&task.ID,
		&task.Name,
		&task.ScriptID,
		&task.UsersCount,
		&task.SpawnRate,
		&task.DurationSeconds,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.StartedAt,
		&task.FinishedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return task, nil
}

func (r *TaskRepo) List(ctx context.Context, limit, offset int) ([]model.Task, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, name, script_id, users_count, spawn_rate, duration_seconds, status, created_at, updated_at, started_at, finished_at FROM locust_tasks ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(
			&task.ID,
			&task.Name,
			&task.ScriptID,
			&task.UsersCount,
			&task.SpawnRate,
			&task.DurationSeconds,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.StartedAt,
			&task.FinishedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

func (r *TaskRepo) Update(ctx context.Context, task *model.Task) error {
	row := r.pool.QueryRow(ctx,
		"UPDATE locust_tasks SET name = $1, script_id = $2, users_count = $3, spawn_rate = $4, duration_seconds = $5, status = $6, started_at = $7, finished_at = $8, updated_at = NOW() WHERE id = $9 RETURNING updated_at",
		task.Name,
		task.ScriptID,
		task.UsersCount,
		task.SpawnRate,
		task.DurationSeconds,
		task.Status,
		task.StartedAt,
		task.FinishedAt,
		task.ID,
	)
	if err := row.Scan(&task.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return repository.ErrNotFound
		}
		return err
	}
	return nil
}

func (r *TaskRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM locust_tasks WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}
