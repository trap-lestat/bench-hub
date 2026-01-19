package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type runRequest struct {
	TaskID          string `json:"task_id"`
	TaskName        string `json:"task_name"`
	UsersCount      int    `json:"users_count"`
	SpawnRate       int    `json:"spawn_rate"`
	DurationSeconds int    `json:"duration_seconds"`
	TargetHost      string `json:"target_host"`
	ScriptContent   string `json:"script_content"`
}

type reportInfo struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	FilePath string `json:"file_path"`
}

type runResponse struct {
	Status  string       `json:"status"`
	Reports []reportInfo `json:"reports"`
}

func main() {
	port := getEnv("RUNNER_PORT", "8081")
	reportsDir := getEnv("REPORTS_DIR", "reports")
	locustBin := getEnv("LOCUST_BIN", "locust")
	locustHost := getEnv("LOCUST_HOST", "http://localhost:8080")

	http.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var req runRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if req.TaskID == "" || req.ScriptContent == "" || req.UsersCount <= 0 || req.SpawnRate <= 0 || req.DurationSeconds <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		reportDir := filepath.Join(reportsDir, fmt.Sprintf("task_%s_%s", req.TaskID, time.Now().Format("20060102150405")))
		if err := os.MkdirAll(reportDir, 0o755); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		scriptPath := filepath.Join(reportDir, "locustfile.py")
		if err := os.WriteFile(scriptPath, []byte(req.ScriptContent), 0o644); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		csvPrefix := filepath.Join(reportDir, "report")
		htmlPath := filepath.Join(reportDir, "report.html")

		targetHost := locustHost
		if req.TargetHost != "" {
			targetHost = req.TargetHost
		}

		cmd := exec.Command(
			locustBin,
			"-f", scriptPath,
			"--headless",
			"-u", fmt.Sprintf("%d", req.UsersCount),
			"-r", fmt.Sprintf("%d", req.SpawnRate),
			"--run-time", fmt.Sprintf("%ds", req.DurationSeconds),
			"--host", targetHost,
			"--csv", csvPrefix,
			"--html", htmlPath,
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		runErr := cmd.Run()

		relativeDir := filepath.Base(reportDir)
		csvFile := filepath.Base(csvPrefix + "_stats.csv")
		htmlFile := filepath.Base(htmlPath)

		reports := make([]reportInfo, 0, 2)
		if _, err := os.Stat(htmlPath); err == nil {
			reports = append(reports, reportInfo{
				Name:     fmt.Sprintf("%s-%s", req.TaskName, htmlFile),
				Type:     "html",
				FilePath: filepath.Join(relativeDir, htmlFile),
			})
		}
		if _, err := os.Stat(csvPrefix + "_stats.csv"); err == nil {
			reports = append(reports, reportInfo{
				Name:     fmt.Sprintf("%s-%s", req.TaskName, csvFile),
				Type:     "csv",
				FilePath: filepath.Join(relativeDir, csvFile),
			})
		}

		resp := runResponse{
			Status:  "finished",
			Reports: reports,
		}
		if runErr != nil {
			resp.Status = "failed"
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})

	log.Printf("runner listening on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("runner stopped: %v", err)
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
