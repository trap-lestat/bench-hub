package api

import (
	"github.com/gin-gonic/gin"

	"bench-hub/internal/api/handlers"
	"bench-hub/internal/middleware"
	"bench-hub/internal/service"
)

func RegisterRoutes(router *gin.Engine, services *service.Services) {
	router.GET("/health", handlers.Health)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", handlers.Ping)

		authHandler := handlers.NewAuthHandler(services.Auth)
		userHandler := handlers.NewUserHandler(services.Users)
		scriptHandler := handlers.NewScriptHandler(services.Scripts)
		taskHandler := handlers.NewTaskHandler(services.Tasks)
		reportHandler := handlers.NewReportHandler(services.Reports)
		taskRunHandler := handlers.NewTaskRunHandler(services.Runner)
		dashboardHandler := handlers.NewDashboardHandler(services.Stats)
		settingsHandler := handlers.NewSettingsHandler(services.Settings)

		v1.POST("/auth/login", authHandler.Login)
		v1.POST("/auth/refresh", authHandler.Refresh)

		protected := v1.Group("")
		protected.Use(middleware.Auth(services.Auth))
		protected.POST("/auth/logout", authHandler.Logout)

		protected.GET("/users", userHandler.List)
		protected.GET("/users/:id", userHandler.Get)
		protected.POST("/users", userHandler.Create)
		protected.PUT("/users/:id", userHandler.Update)
		protected.DELETE("/users/:id", userHandler.Delete)

		protected.GET("/scripts", scriptHandler.List)
		protected.GET("/scripts/:id", scriptHandler.Get)
		protected.POST("/scripts", scriptHandler.Create)
		protected.PUT("/scripts/:id", scriptHandler.Update)
		protected.DELETE("/scripts/:id", scriptHandler.Delete)
		protected.POST("/scripts/import", scriptHandler.Import)

		protected.GET("/tasks", taskHandler.List)
		protected.GET("/tasks/:id", taskHandler.Get)
		protected.POST("/tasks", taskHandler.Create)
		protected.POST("/tasks/:id/stop", taskHandler.Stop)
		protected.POST("/tasks/:id/run", taskRunHandler.Run)

		protected.GET("/reports", reportHandler.List)
		protected.GET("/reports/:id", reportHandler.Get)
		protected.GET("/reports/:id/download", reportHandler.Download)
		protected.GET("/reports/:id/preview", reportHandler.Preview)

		protected.GET("/dashboard/summary", dashboardHandler.Summary)
		protected.GET("/settings/p95-baseline", settingsHandler.GetP95)
		protected.PUT("/settings/p95-baseline", settingsHandler.UpdateP95)
	}
}
