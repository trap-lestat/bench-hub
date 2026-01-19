package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"bench-hub/internal/api"
	"bench-hub/internal/config"
	"bench-hub/internal/middleware"
	"bench-hub/internal/migrate"
	"bench-hub/internal/observability"
	"bench-hub/internal/repository/postgres"
	"bench-hub/internal/service"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, cfg)
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}
	defer pool.Close()

	if cfg.AutoMigrate {
		if err := migrate.Run(ctx, cfg); err != nil {
			log.Fatalf("auto migrate failed: %v", err)
		}
	}

	userRepo := postgres.NewUserRepo(pool)
	scriptRepo := postgres.NewScriptRepo(pool)
	taskRepo := postgres.NewTaskRepo(pool)
	reportRepo := postgres.NewReportRepo(pool)
	settingsRepo := postgres.NewSettingsRepo(pool)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.AccessTokenMinutes, cfg.RefreshTokenDays, cfg.JWTIssuer)
	userService := service.NewUserService(userRepo)
	scriptService := service.NewScriptService(scriptRepo)
	taskService := service.NewTaskService(taskRepo)
	reportService := service.NewReportService(reportRepo, cfg.ReportsDir)
	runner := service.NewTaskRunner(taskRepo, scriptRepo, reportRepo, cfg.ReportsDir, cfg.LocustBin, cfg.LocustHost, cfg.RunnerURL)
	settingsService := service.NewSettingsService(settingsRepo)
	statsService := service.NewStatsService(userRepo, scriptRepo, reportRepo, settingsService)

	services := &service.Services{
		Auth:     authService,
		Users:    userService,
		Scripts:  scriptService,
		Tasks:    taskService,
		Reports:  reportService,
		Runner:   runner,
		Settings: settingsService,
		Stats:    statsService,
	}

	router := gin.New()
	router.Use(middleware.RequestID())
	router.Use(middleware.Recovery())
	router.Use(gin.Logger())
	router.Use(observability.Metrics())

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	api.RegisterRoutes(router, services)

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
