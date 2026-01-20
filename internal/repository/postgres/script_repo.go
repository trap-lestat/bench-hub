package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"bench-hub/internal/model"
	"bench-hub/internal/repository"
)

type ScriptRepo struct {
	pool *pgxpool.Pool
}

func NewScriptRepo(pool *pgxpool.Pool) *ScriptRepo {
	return &ScriptRepo{pool: pool}
}

func (r *ScriptRepo) Create(ctx context.Context, script *model.Script) error {
	if script.ID == "" {
		script.ID = uuid.NewString()
	}

	row := r.pool.QueryRow(ctx,
		"INSERT INTO locust_scripts (id, name, description, script_type, content) VALUES ($1, $2, $3, $4, $5) RETURNING created_at, updated_at",
		script.ID,
		script.Name,
		script.Description,
		script.Type,
		script.Content,
	)

	return row.Scan(&script.CreatedAt, &script.UpdatedAt)
}

func (r *ScriptRepo) GetByID(ctx context.Context, id string) (*model.Script, error) {
	script := &model.Script{}
	row := r.pool.QueryRow(ctx,
		"SELECT id, name, description, script_type, content, created_at, updated_at FROM locust_scripts WHERE id = $1",
		id,
	)
	if err := row.Scan(&script.ID, &script.Name, &script.Description, &script.Type, &script.Content, &script.CreatedAt, &script.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return script, nil
}

func (r *ScriptRepo) List(ctx context.Context, limit, offset int) ([]model.Script, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, name, description, script_type, content, created_at, updated_at FROM locust_scripts ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scripts []model.Script
	for rows.Next() {
		var script model.Script
		if err := rows.Scan(&script.ID, &script.Name, &script.Description, &script.Type, &script.Content, &script.CreatedAt, &script.UpdatedAt); err != nil {
			return nil, err
		}
		scripts = append(scripts, script)
	}
	return scripts, rows.Err()
}

func (r *ScriptRepo) Update(ctx context.Context, script *model.Script) error {
	row := r.pool.QueryRow(ctx,
		"UPDATE locust_scripts SET name = $1, description = $2, script_type = $3, content = $4, updated_at = NOW() WHERE id = $5 RETURNING updated_at",
		script.Name,
		script.Description,
		script.Type,
		script.Content,
		script.ID,
	)
	if err := row.Scan(&script.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return repository.ErrNotFound
		}
		return err
	}
	return nil
}

func (r *ScriptRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM locust_scripts WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *ScriptRepo) Count(ctx context.Context) (int, error) {
	var count int
	row := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM locust_scripts")
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
