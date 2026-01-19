package service

import (
	"context"

	"bench-hub/internal/repository"
)

type Summary struct {
	UsersCount   int    `json:"users_count"`
	ScriptsCount int    `json:"scripts_count"`
	ReportsCount int    `json:"reports_count"`
	P95Baseline  string `json:"p95_baseline"`
}

type StatsService struct {
	users    repository.UserRepository
	scripts  repository.ScriptRepository
	reports  repository.ReportRepository
	settings *SettingsService
}

func NewStatsService(users repository.UserRepository, scripts repository.ScriptRepository, reports repository.ReportRepository, settings *SettingsService) *StatsService {
	return &StatsService{
		users:    users,
		scripts:  scripts,
		reports:  reports,
		settings: settings,
	}
}

func (s *StatsService) Summary(ctx context.Context) (Summary, error) {
	usersCount, err := s.users.Count(ctx)
	if err != nil {
		return Summary{}, err
	}
	scriptsCount, err := s.scripts.Count(ctx)
	if err != nil {
		return Summary{}, err
	}
	reportsCount, err := s.reports.Count(ctx)
	if err != nil {
		return Summary{}, err
	}
	p95, err := s.settings.GetP95Baseline(ctx)
	if err != nil {
		return Summary{}, err
	}

	return Summary{
		UsersCount:   usersCount,
		ScriptsCount: scriptsCount,
		ReportsCount: reportsCount,
		P95Baseline:  p95,
	}, nil
}
