package service

import (
	"context"
	"path/filepath"
	"strings"

	"bench-hub/internal/model"
	"bench-hub/internal/repository"
)

type ReportService struct {
	repo       repository.ReportRepository
	reportsDir string
}

func NewReportService(repo repository.ReportRepository, reportsDir string) *ReportService {
	return &ReportService{repo: repo, reportsDir: reportsDir}
}

func (s *ReportService) Create(ctx context.Context, taskID *string, name, reportType, filePath string) (*model.Report, error) {
	report := &model.Report{
		TaskID:   taskID,
		Name:     name,
		Type:     reportType,
		FilePath: filePath,
	}

	if err := s.repo.Create(ctx, report); err != nil {
		return nil, err
	}

	return report, nil
}

func (s *ReportService) Get(ctx context.Context, id string) (*model.Report, error) {
	report, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return report, nil
}

func (s *ReportService) List(ctx context.Context, limit, offset int) ([]model.Report, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *ReportService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if err == repository.ErrNotFound {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *ReportService) ResolvePath(filePath string) (string, bool) {
	clean := filepath.Clean(filePath)
	if filepath.IsAbs(clean) {
		clean = strings.TrimPrefix(clean, string(filepath.Separator))
	}
	full := filepath.Join(s.reportsDir, clean)

	base := filepath.Clean(s.reportsDir) + string(filepath.Separator)
	fullClean := filepath.Clean(full)
	if !strings.HasPrefix(fullClean+string(filepath.Separator), base) {
		return "", false
	}
	return fullClean, true
}
