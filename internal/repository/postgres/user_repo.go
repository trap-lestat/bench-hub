package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"bench-hub/internal/model"
	"bench-hub/internal/repository"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	if user.ID == "" {
		user.ID = uuid.NewString()
	}

	row := r.pool.QueryRow(ctx,
		"INSERT INTO users (id, username, password_hash) VALUES ($1, $2, $3) RETURNING created_at",
		user.ID,
		user.Username,
		user.PasswordHash,
	)

	return row.Scan(&user.CreatedAt)
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	user := &model.User{}
	row := r.pool.QueryRow(ctx,
		"SELECT id, username, password_hash, created_at FROM users WHERE id = $1",
		id,
	)
	if err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user := &model.User{}
	row := r.pool.QueryRow(ctx,
		"SELECT id, username, password_hash, created_at FROM users WHERE username = $1",
		username,
	)
	if err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) List(ctx context.Context, limit, offset int) ([]model.User, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, username, password_hash, created_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	tag, err := r.pool.Exec(ctx,
		"UPDATE users SET username = $1, password_hash = $2 WHERE id = $3",
		user.Username,
		user.PasswordHash,
		user.ID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *UserRepo) Count(ctx context.Context) (int, error) {
	var count int
	row := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM users")
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
