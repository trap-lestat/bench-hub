package repository

import (
	"context"

	"bench-hub/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	List(ctx context.Context, limit, offset int) ([]model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int, error)
}

type ScriptRepository interface {
	Create(ctx context.Context, script *model.Script) error
	GetByID(ctx context.Context, id string) (*model.Script, error)
	List(ctx context.Context, limit, offset int) ([]model.Script, error)
	Update(ctx context.Context, script *model.Script) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int, error)
}

type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
	GetByID(ctx context.Context, id string) (*model.Task, error)
	List(ctx context.Context, limit, offset int) ([]model.Task, error)
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, id string) error
}

type ReportRepository interface {
	Create(ctx context.Context, report *model.Report) error
	GetByID(ctx context.Context, id string) (*model.Report, error)
	List(ctx context.Context, limit, offset int) ([]model.Report, error)
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int, error)
}

type SettingsRepository interface {
	Get(ctx context.Context, key string) (string, bool, error)
	Set(ctx context.Context, key, value string) error
}
