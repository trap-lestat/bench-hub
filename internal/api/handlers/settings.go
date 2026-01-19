package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"bench-hub/internal/model"
	"bench-hub/internal/service"
)

type SettingsHandler struct {
	settings *service.SettingsService
}

type updateP95Request struct {
	Value string `json:"value"`
}

func NewSettingsHandler(settings *service.SettingsService) *SettingsHandler {
	return &SettingsHandler{settings: settings}
}

func (h *SettingsHandler) GetP95(c *gin.Context) {
	value, err := h.settings.GetP95Baseline(c.Request.Context())
	if err != nil {
		model.JSON(c, http.StatusInternalServerError, model.Fail(9000, "internal error"))
		return
	}
	model.JSON(c, http.StatusOK, model.OK(gin.H{"value": value}))
}

func (h *SettingsHandler) UpdateP95(c *gin.Context) {
	var req updateP95Request
	if err := c.ShouldBindJSON(&req); err != nil {
		model.JSON(c, http.StatusBadRequest, model.Fail(1000, "invalid params"))
		return
	}
	if err := h.settings.SetP95Baseline(c.Request.Context(), req.Value); err != nil {
		model.JSON(c, http.StatusInternalServerError, model.Fail(9000, "internal error"))
		return
	}
	model.JSON(c, http.StatusOK, model.OK(nil))
}
