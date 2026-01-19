package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port               string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPass             string
	DBName             string
	DBSSLMode          string
	JWTSecret          string
	AccessTokenMinutes time.Duration
	RefreshTokenDays   time.Duration
	JWTIssuer          string
	ReportsDir         string
	LocustBin          string
	LocustHost         string
	MigrationsPath     string
	AutoMigrate        bool
	RunnerURL          string
}

func Load() Config {
	return Config{
		Port:               getEnv("PORT", "8080"),
		DBHost:             getEnv("DB_HOST", "127.0.0.1"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPass:             getEnv("DB_PASS", "postgres"),
		DBName:             getEnv("DB_NAME", "bench_hub"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		JWTSecret:          getEnv("JWT_SECRET", "dev-secret"),
		AccessTokenMinutes: time.Duration(getEnvInt("ACCESS_TOKEN_MINUTES", 60)) * time.Minute,
		RefreshTokenDays:   time.Duration(getEnvInt("REFRESH_TOKEN_DAYS", 7)) * 24 * time.Hour,
		JWTIssuer:          getEnv("JWT_ISSUER", "bench-hub"),
		ReportsDir:         getEnv("REPORTS_DIR", "reports"),
		LocustBin:          getEnv("LOCUST_BIN", "locust"),
		LocustHost:         getEnv("LOCUST_HOST", "http://localhost:8080"),
		MigrationsPath:     getEnv("MIGRATIONS_PATH", "migrations"),
		AutoMigrate:        getEnvBool("AUTO_MIGRATE", false),
		RunnerURL:          getEnv("RUNNER_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return parsed
}

func getEnvBool(key string, fallback bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	switch val {
	case "1", "true", "TRUE", "yes", "YES":
		return true
	case "0", "false", "FALSE", "no", "NO":
		return false
	default:
		return fallback
	}
}
