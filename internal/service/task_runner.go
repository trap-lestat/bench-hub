package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"bench-hub/internal/model"
	"bench-hub/internal/repository"
)

type TaskRunner struct {
	tasks      repository.TaskRepository
	scripts    repository.ScriptRepository
	reports    repository.ReportRepository
	reportsDir string
	locustBin  string
	locustHost string
	runnerURL  string
}

func NewTaskRunner(tasks repository.TaskRepository, scripts repository.ScriptRepository, reports repository.ReportRepository, reportsDir, locustBin, locustHost, runnerURL string) *TaskRunner {
	return &TaskRunner{
		tasks:      tasks,
		scripts:    scripts,
		reports:    reports,
		reportsDir: reportsDir,
		locustBin:  locustBin,
		locustHost: locustHost,
		runnerURL:  runnerURL,
	}
}

func (r *TaskRunner) Run(ctx context.Context, taskID, targetHost string) (*model.Task, error) {
	task, err := r.tasks.GetByID(ctx, taskID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if task.Status == TaskStatusRunning {
		return task, nil
	}

	script, err := r.scripts.GetByID(ctx, task.ScriptID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	now := time.Now()
	task.Status = TaskStatusRunning
	task.StartedAt = &now
	task.FinishedAt = nil
	if err := r.tasks.Update(ctx, task); err != nil {
		return nil, err
	}

	go r.execute(task, script, targetHost)

	return task, nil
}

func (r *TaskRunner) execute(task *model.Task, script *model.Script, targetHost string) {
	runCtx := context.Background()
	finishTime := time.Now()
	status := TaskStatusFinished

	if r.runnerURL != "" {
		reports, runStatus, err := r.runRemote(task, script, targetHost)
		if err != nil {
			status = TaskStatusFailed
		} else {
			if runStatus == "failed" {
				status = TaskStatusFailed
			}
			for _, report := range reports {
				_ = r.reports.Create(runCtx, &model.Report{
					TaskID:   &task.ID,
					Name:     report.Name,
					Type:     report.Type,
					FilePath: report.FilePath,
				})
			}
		}
	} else if err := r.runLocust(task, script, targetHost); err != nil {
		status = TaskStatusFailed
	}

	task.Status = status
	task.FinishedAt = &finishTime
	_ = r.tasks.Update(runCtx, task)
}

type runnerRequest struct {
	TaskID          string `json:"task_id"`
	TaskName        string `json:"task_name"`
	UsersCount      int    `json:"users_count"`
	SpawnRate       int    `json:"spawn_rate"`
	DurationSeconds int    `json:"duration_seconds"`
	TargetHost      string `json:"target_host"`
	ScriptContent   string `json:"script_content"`
}

type runnerReport struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	FilePath string `json:"file_path"`
}

type runnerResponse struct {
	Status  string         `json:"status"`
	Reports []runnerReport `json:"reports"`
}

func (r *TaskRunner) runRemote(task *model.Task, script *model.Script, targetHost string) ([]runnerReport, string, error) {
	client := &http.Client{Timeout: time.Duration(task.DurationSeconds+30) * time.Second}
	reqBody := runnerRequest{
		TaskID:          task.ID,
		TaskName:        task.Name,
		UsersCount:      task.UsersCount,
		SpawnRate:       task.SpawnRate,
		DurationSeconds: task.DurationSeconds,
		TargetHost:      targetHost,
		ScriptContent:   script.Content,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, "", err
	}

	resp, err := client.Post(r.runnerURL+"/run", "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("runner status %d", resp.StatusCode)
	}

	var out runnerResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, "", err
	}

	return out.Reports, out.Status, nil
}

func (r *TaskRunner) runLocust(task *model.Task, script *model.Script, targetHost string) error {
	timestamp := time.Now().Format("20060102150405")
	reportDir := filepath.Join(r.reportsDir, fmt.Sprintf("task_%s_%s", task.ID, timestamp))
	if err := os.MkdirAll(reportDir, 0o755); err != nil {
		return err
	}

	scriptPath := filepath.Join(reportDir, "locustfile.py")
	if err := os.WriteFile(scriptPath, []byte(script.Content), 0o644); err != nil {
		return err
	}

	csvPrefix := filepath.Join(reportDir, "report")
	htmlPath := filepath.Join(reportDir, "report.html")

	host := r.locustHost
	if targetHost != "" {
		host = targetHost
	}

	cmd := exec.Command(
		r.locustBin,
		"-f", scriptPath,
		"--headless",
		"-u", fmt.Sprintf("%d", task.UsersCount),
		"-r", fmt.Sprintf("%d", task.SpawnRate),
		"--run-time", fmt.Sprintf("%ds", task.DurationSeconds),
		"--host", host,
		"--csv", csvPrefix,
		"--html", htmlPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	csvStats := csvPrefix + "_stats.csv"
	csvFile := filepath.Base(csvStats)
	htmlFile := filepath.Base(htmlPath)

	relativeDir := filepath.Base(reportDir)
	_ = r.reports.Create(context.Background(), &model.Report{
		TaskID:   &task.ID,
		Name:     fmt.Sprintf("%s-%s", task.Name, htmlFile),
		Type:     "html",
		FilePath: filepath.Join(relativeDir, htmlFile),
	})
	_ = r.reports.Create(context.Background(), &model.Report{
		TaskID:   &task.ID,
		Name:     fmt.Sprintf("%s-%s", task.Name, csvFile),
		Type:     "csv",
		FilePath: filepath.Join(relativeDir, csvFile),
	})

	return nil
}
