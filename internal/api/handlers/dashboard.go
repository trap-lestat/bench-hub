package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"bench-hub/internal/model"
	"bench-hub/internal/service"
)

type DashboardHandler struct {
	stats *service.StatsService
}

func NewDashboardHandler(stats *service.StatsService) *DashboardHandler {
	return &DashboardHandler{stats: stats}
}

func (h *DashboardHandler) Summary(c *gin.Context) {
	summary, err := h.stats.Summary(c.Request.Context())
	if err != nil {
		model.JSON(c, http.StatusInternalServerError, model.Fail(9000, "internal error"))
		return
	}
	model.JSON(c, http.StatusOK, model.OK(summary))
}
