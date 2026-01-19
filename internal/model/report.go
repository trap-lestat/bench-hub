package model

import "time"

type Report struct {
	ID        string    `json:"id"`
	TaskID    *string   `json:"task_id"`
	TaskName  *string   `json:"task_name"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	FilePath  string    `json:"file_path"`
	CreatedAt time.Time `json:"created_at"`
}
