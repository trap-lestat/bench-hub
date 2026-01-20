package service

import (
	"context"
	"strings"

	"bench-hub/internal/model"
	"bench-hub/internal/repository"
)

type ScriptService struct {
	repo repository.ScriptRepository
}

func NewScriptService(repo repository.ScriptRepository) *ScriptService {
	return &ScriptService{repo: repo}
}

func normalizeScriptType(value string) (string, error) {
	if value == "" {
		return model.ScriptTypeLocust, nil
	}

	value = strings.ToLower(strings.TrimSpace(value))
	switch value {
	case model.ScriptTypeLocust, model.ScriptTypeJMeter:
		return value, nil
	default:
		return "", ErrInvalidScriptType
	}
}

func (s *ScriptService) Create(ctx context.Context, name, description, scriptType, content string) (*model.Script, error) {
	kind, err := normalizeScriptType(scriptType)
	if err != nil {
		return nil, err
	}

	script := &model.Script{
		Name:        name,
		Description: description,
		Type:        kind,
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

func (s *ScriptService) Update(ctx context.Context, id, name, description, scriptType, content string) (*model.Script, error) {
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
	if scriptType != "" {
		kind, err := normalizeScriptType(scriptType)
		if err != nil {
			return nil, err
		}
		script.Type = kind
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
