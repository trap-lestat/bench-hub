package service

import (
	"context"

	"bench-hub/internal/repository"
)

const settingsKeyP95Baseline = "p95_baseline"

// Default is used when not configured.
const defaultP95Baseline = "P95 < 300ms"

type SettingsService struct {
	repo repository.SettingsRepository
}

func NewSettingsService(repo repository.SettingsRepository) *SettingsService {
	return &SettingsService{repo: repo}
}

func (s *SettingsService) GetP95Baseline(ctx context.Context) (string, error) {
	value, ok, err := s.repo.Get(ctx, settingsKeyP95Baseline)
	if err != nil {
		return "", err
	}
	if !ok || value == "" {
		return defaultP95Baseline, nil
	}
	return value, nil
}

func (s *SettingsService) SetP95Baseline(ctx context.Context, value string) error {
	if value == "" {
		value = defaultP95Baseline
	}
	return s.repo.Set(ctx, settingsKeyP95Baseline, value)
}
