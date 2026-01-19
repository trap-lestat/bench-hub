package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"bench-hub/internal/model"
	"bench-hub/internal/repository"
)

type ReportRepo struct {
	pool *pgxpool.Pool
}

func NewReportRepo(pool *pgxpool.Pool) *ReportRepo {
	return &ReportRepo{pool: pool}
}

func (r *ReportRepo) Create(ctx context.Context, report *model.Report) error {
	if report.ID == "" {
		report.ID = uuid.NewString()
	}

	row := r.pool.QueryRow(ctx,
		"INSERT INTO locust_reports (id, task_id, name, report_type, file_path) VALUES ($1, $2, $3, $4, $5) RETURNING created_at",
		report.ID,
		report.TaskID,
		report.Name,
		report.Type,
		report.FilePath,
	)

	return row.Scan(&report.CreatedAt)
}

func (r *ReportRepo) GetByID(ctx context.Context, id string) (*model.Report, error) {
	report := &model.Report{}
	row := r.pool.QueryRow(ctx,
		`SELECT r.id, r.task_id, t.name, r.name, r.report_type, r.file_path, r.created_at
		 FROM locust_reports r
		 LEFT JOIN locust_tasks t ON r.task_id = t.id
		 WHERE r.id = $1`,
		id,
	)
	if err := row.Scan(&report.ID, &report.TaskID, &report.TaskName, &report.Name, &report.Type, &report.FilePath, &report.CreatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return report, nil
}

func (r *ReportRepo) List(ctx context.Context, limit, offset int) ([]model.Report, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT r.id, r.task_id, t.name, r.name, r.report_type, r.file_path, r.created_at
		 FROM locust_reports r
		 LEFT JOIN locust_tasks t ON r.task_id = t.id
		 ORDER BY r.created_at DESC
		 LIMIT $1 OFFSET $2`,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []model.Report
	for rows.Next() {
		var report model.Report
		if err := rows.Scan(&report.ID, &report.TaskID, &report.TaskName, &report.Name, &report.Type, &report.FilePath, &report.CreatedAt); err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	return reports, rows.Err()
}

func (r *ReportRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM locust_reports WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *ReportRepo) Count(ctx context.Context) (int, error) {
	var count int
	row := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM locust_reports")
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
