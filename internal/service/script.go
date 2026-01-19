package service

import (
	"context"

	"bench-hub/internal/model"
	"bench-hub/internal/repository"
)

type ScriptService struct {
	repo repository.ScriptRepository
}

func NewScriptService(repo repository.ScriptRepository) *ScriptService {
	return &ScriptService{repo: repo}
}

func (s *ScriptService) Create(ctx context.Context, name, description, content string) (*model.Script, error) {
	script := &model.Script{
		Name:        name,
		Description: description,
		Content:     content,
	}

	if err := s.repo.Create(ctx, script); err != nil {
		return nil, err
	}

	return script, nil
}

func (s *ScriptService) Get(ctx context.Context, id string) (*model.Script, error) {
	script, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return script, nil
}

func (s *ScriptService) List(ctx context.Context, limit, offset int) ([]model.Script, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *ScriptService) Update(ctx context.Context, id, name, description, content string) (*model.Script, error) {
	script, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if name != "" {
		script.Name = name
	}
	if description != "" {
		script.Description = description
	}
	if content != "" {
		script.Content = content
	}

	if err := s.repo.Update(ctx, script); err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return script, nil
}

func (s *ScriptService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if err == repository.ErrNotFound {
			return ErrNotFound
		}
		return err
	}
	return nil
}
