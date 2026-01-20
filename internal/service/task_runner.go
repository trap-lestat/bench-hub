package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
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
	runningMu  sync.Mutex
	running    map[string]*runningCommand
}

type runningCommand struct {
	cmd     *exec.Cmd
	stopped bool
}

func pickTargetHost(input string, fallback *string) string {
	if strings.TrimSpace(input) != "" {
		return strings.TrimSpace(input)
	}
	if fallback != nil {
		return strings.TrimSpace(*fallback)
	}
	return ""
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
		running:    make(map[string]*runningCommand),
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

	go r.execute(task, script, pickTargetHost(targetHost, task.TargetHost))

	return task, nil
}

func (r *TaskRunner) Stop(ctx context.Context, taskID string) (*model.Task, error) {
	task, err := r.tasks.GetByID(ctx, taskID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if task.Status == TaskStatusFinished || task.Status == TaskStatusFailed || task.Status == TaskStatusStopped {
		return task, nil
	}

	if r.runnerURL != "" {
		if err := r.stopRemote(taskID); err != nil {
			return nil, err
		}
	} else {
		_ = r.stopLocal(taskID)
	}

	now := time.Now()
	task.Status = TaskStatusStopped
	if task.StartedAt == nil {
		task.StartedAt = &now
	}
	task.FinishedAt = &now
	if err := r.tasks.Update(ctx, task); err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return task, nil
}

func (r *TaskRunner) setRunning(taskID string, cmd *exec.Cmd) {
	r.runningMu.Lock()
	r.running[taskID] = &runningCommand{cmd: cmd}
	r.runningMu.Unlock()
}

func (r *TaskRunner) clearRunning(taskID string) (stopped bool) {
	r.runningMu.Lock()
	entry := r.running[taskID]
	if entry != nil {
		stopped = entry.stopped
		delete(r.running, taskID)
	}
	r.runningMu.Unlock()
	return stopped
}

func (r *TaskRunner) markStopped(taskID string) *exec.Cmd {
	r.runningMu.Lock()
	entry := r.running[taskID]
	var cmd *exec.Cmd
	if entry != nil {
		entry.stopped = true
		cmd = entry.cmd
	}
	r.runningMu.Unlock()
	return cmd
}

func (r *TaskRunner) stopLocal(taskID string) error {
	cmd := r.markStopped(taskID)
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		_ = cmd.Process.Kill()
	}
	return nil
}

func (r *TaskRunner) stopRemote(taskID string) error {
	client := &http.Client{Timeout: 10 * time.Second}
	body, err := json.Marshal(map[string]string{"task_id": taskID})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, r.runnerURL+"/stop", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("runner stop status %d", resp.StatusCode)
	}
	return nil
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
			} else if runStatus == "stopped" {
				status = TaskStatusStopped
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
	} else if err := r.runLocal(task, script, targetHost); err != nil {
		if errors.Is(err, ErrStopped) {
			status = TaskStatusStopped
		} else {
			status = TaskStatusFailed
		}
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
	JmeterTPM       *int   `json:"jmeter_tpm"`
	ScriptType      string `json:"script_type"`
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
		JmeterTPM:       task.JmeterTPM,
		ScriptType:      script.Type,
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

func (r *TaskRunner) runLocal(task *model.Task, script *model.Script, targetHost string) error {
	if script.Type == "" || script.Type == model.ScriptTypeLocust {
		return r.runLocust(task, script, targetHost)
	}
	if script.Type == model.ScriptTypeJMeter {
		return ErrUnsupportedEngine
	}
	return ErrInvalidScriptType
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

	r.setRunning(task.ID, cmd)
	cmdErr := cmd.Run()
	stopped := r.clearRunning(task.ID)

	if cmdErr != nil && !stopped {
		return cmdErr
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

	if stopped {
		return ErrStopped
	}
	return nil
}
