package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"bench-hub/internal/model"
	"bench-hub/internal/service"
)

type TaskRunHandler struct {
	runner *service.TaskRunner
}

func NewTaskRunHandler(runner *service.TaskRunner) *TaskRunHandler {
	return &TaskRunHandler{runner: runner}
}

type runTaskRequest struct {
	TargetHost string `json:"target_host"`
}

func (h *TaskRunHandler) Run(c *gin.Context) {
	id := c.Param("id")
	var req runTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		model.JSON(c, http.StatusBadRequest, model.Fail(1000, "invalid params"))
		return
	}
	task, err := h.runner.Run(c.Request.Context(), id, req.TargetHost)
	if err != nil {
		if err == service.ErrNotFound {
			model.JSON(c, http.StatusNotFound, model.Fail(1003, "not found"))
			return
		}
		model.JSON(c, http.StatusInternalServerError, model.Fail(9000, "internal error"))
		return
	}

	model.JSON(c, http.StatusOK, model.OK(task))
}
