package model

import "time"

type Task struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	ScriptID        string     `json:"script_id"`
	UsersCount      int        `json:"users_count"`
	SpawnRate       int        `json:"spawn_rate"`
	DurationSeconds int        `json:"duration_seconds"`
	TargetHost      *string    `json:"target_host"`
	JmeterTPM       *int       `json:"jmeter_tpm"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	StartedAt       *time.Time `json:"started_at"`
	FinishedAt      *time.Time `json:"finished_at"`
}
