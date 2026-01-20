package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type runRequest struct {
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

type reportInfo struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	FilePath string `json:"file_path"`
}

type runResponse struct {
	Status  string       `json:"status"`
	Reports []reportInfo `json:"reports"`
}

type runningTask struct {
	cmd     *exec.Cmd
	stopped bool
}

func main() {
	port := getEnv("RUNNER_PORT", "8081")
	reportsDir := getEnv("REPORTS_DIR", "reports")
	locustBin := getEnv("LOCUST_BIN", "locust")
	jmeterBin := getEnv("JMETER_BIN", "jmeter")
	locustHost := getEnv("LOCUST_HOST", "http://localhost:8080")

	runningMu := sync.Mutex{}
	running := map[string]*runningTask{}

	markStopped := func(taskID string) *exec.Cmd {
		runningMu.Lock()
		entry := running[taskID]
		var cmd *exec.Cmd
		if entry != nil {
			entry.stopped = true
			cmd = entry.cmd
		}
		runningMu.Unlock()
		return cmd
	}

	register := func(taskID string, cmd *exec.Cmd) {
		runningMu.Lock()
		running[taskID] = &runningTask{cmd: cmd}
		runningMu.Unlock()
	}

	clear := func(taskID string) bool {
		runningMu.Lock()
		entry := running[taskID]
		stopped := false
		if entry != nil {
			stopped = entry.stopped
			delete(running, taskID)
		}
		runningMu.Unlock()
		return stopped
	}

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

		scriptType, ok := normalizeScriptType(req.ScriptType)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		reportDir := filepath.Join(reportsDir, fmt.Sprintf("task_%s_%s", req.TaskID, time.Now().Format("20060102150405")))
		if err := os.MkdirAll(reportDir, 0o755); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var cmd *exec.Cmd
		var reports []reportInfo

		if scriptType == "jmeter" {
			scriptPath := filepath.Join(reportDir, "test.jmx")
			if err := os.WriteFile(scriptPath, []byte(req.ScriptContent), 0o644); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			resultsPath := filepath.Join(reportDir, "results.jtl")
			htmlDir := filepath.Join(reportDir, "html-report")
			host, port, protocol := parseTargetHost(req.TargetHost, locustHost)
			args := []string{"-n", "-t", scriptPath, "-l", resultsPath, "-e", "-o", htmlDir, "-Jtarget_host=" + host, "-Jtarget_port=" + port, "-Jtarget_protocol=" + protocol, "-Jduration=" + fmt.Sprintf("%d", req.DurationSeconds)}
			if req.JmeterTPM != nil && *req.JmeterTPM > 0 {
				args = append(args, "-Jtpm="+fmt.Sprintf("%d", *req.JmeterTPM))
			}

			cmd = exec.Command(jmeterBin, args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		} else {
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

			cmd = exec.Command(
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
		}

		register(req.TaskID, cmd)
		runErr := cmd.Run()
		stopped := clear(req.TaskID)

		relativeDir := filepath.Base(reportDir)
		var jmeterChecked bool
		var jmeterFailed bool
		if scriptType == "jmeter" {
			htmlPath := filepath.Join(reportDir, "html-report", "index.html")
			if _, err := os.Stat(htmlPath); err == nil {
				reports = append(reports, reportInfo{
					Name:     fmt.Sprintf("%s-jmeter-report.html", req.TaskName),
					Type:     "html",
					FilePath: filepath.Join(relativeDir, "html-report", "index.html"),
				})
			}
			resultsPath := filepath.Join(reportDir, "results.jtl")
			if _, err := os.Stat(resultsPath); err == nil {
				reports = append(reports, reportInfo{
					Name:     fmt.Sprintf("%s-results.jtl", req.TaskName),
					Type:     "jtl",
					FilePath: filepath.Join(relativeDir, "results.jtl"),
				})
				jmeterChecked = true
				if failed, err := jmeterHasFailures(resultsPath); err != nil {
					log.Printf("jmeter jtl parse error: %v", err)
				} else {
					jmeterFailed = failed
				}
			}
		} else {
			csvPrefix := filepath.Join(reportDir, "report")
			htmlPath := filepath.Join(reportDir, "report.html")
			csvFile := filepath.Base(csvPrefix + "_stats.csv")
			htmlFile := filepath.Base(htmlPath)

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
		}

		resp := runResponse{
			Status:  "finished",
			Reports: reports,
		}
		if stopped {
			resp.Status = "stopped"
		} else if scriptType == "jmeter" && jmeterChecked {
			if jmeterFailed {
				resp.Status = "failed"
			}
		} else if runErr != nil {
			resp.Status = "failed"
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			TaskID string `json:"task_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.TaskID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		cmd := markStopped(req.TaskID)
		if cmd == nil || cmd.Process == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			_ = cmd.Process.Kill()
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "stopped"})
	})

	log.Printf("runner listening on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("runner stopped: %v", err)
	}
}

func parseTargetHost(input, fallback string) (string, string, string) {
	target := strings.TrimSpace(input)
	if target == "" {
		target = strings.TrimSpace(fallback)
	}

	protocol := "http"
	hostPort := target

	if strings.Contains(target, "://") {
		if parsed, err := url.Parse(target); err == nil {
			if parsed.Scheme != "" {
				protocol = parsed.Scheme
			}
			if parsed.Host != "" {
				hostPort = parsed.Host
			} else if parsed.Path != "" {
				hostPort = parsed.Path
			}
		}
	}

	if strings.Contains(hostPort, "/") {
		hostPort = strings.Split(hostPort, "/")[0]
	}

	host := hostPort
	port := ""
	if h, p, err := net.SplitHostPort(hostPort); err == nil {
		host = h
		port = p
	}

	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		host = strings.TrimPrefix(strings.TrimSuffix(host, "]"), "[")
	}

	if port == "" {
		if protocol == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}

	if _, err := strconv.Atoi(port); err != nil {
		port = "80"
	}

	return host, port, protocol
}

func jmeterHasFailures(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.FieldsPerRecord = -1

	header, err := reader.Read()
	if err != nil {
		return false, err
	}

	successIndex := -1
	for i, field := range header {
		if strings.EqualFold(strings.TrimSpace(field), "success") {
			successIndex = i
			break
		}
	}
	if successIndex == -1 {
		return false, fmt.Errorf("success column not found in jtl")
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return false, err
		}
		if successIndex >= len(record) {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(record[successIndex]), "false") {
			return true, nil
		}
	}

	return false, nil
}

func normalizeScriptType(value string) (string, bool) {
	if value == "" {
		return "locust", true
	}
	value = strings.ToLower(strings.TrimSpace(value))
	switch value {
	case "locust", "jmeter":
		return value, true
	default:
		return "", false
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
