package service

import (
	"context"
	"time"

	"bench-hub/internal/model"
	"bench-hub/internal/repository"
)

const (
	TaskStatusCreated  = "created"
	TaskStatusRunning  = "running"
	TaskStatusStopped  = "stopped"
	TaskStatusFinished = "finished"
	TaskStatusFailed   = "failed"
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) Create(ctx context.Context, name, scriptID string, usersCount, spawnRate, durationSeconds int) (*model.Task, error) {
	task := &model.Task{
		Name:            name,
		ScriptID:        scriptID,
		UsersCount:      usersCount,
		SpawnRate:       spawnRate,
		DurationSeconds: durationSeconds,
		Status:          TaskStatusCreated,
	}

	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) Get(ctx context.Context, id string) (*model.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return task, nil
}

func (s *TaskService) List(ctx context.Context, limit, offset int) ([]model.Task, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *TaskService) Stop(ctx context.Context, id string) (*model.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if task.Status == TaskStatusFinished || task.Status == TaskStatusFailed {
		return task, nil
	}

	now := time.Now()
	task.Status = TaskStatusStopped
	if task.StartedAt == nil {
		task.StartedAt = &now
	}
	task.FinishedAt = &now

	if err := s.repo.Update(ctx, task); err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return task, nil
}
