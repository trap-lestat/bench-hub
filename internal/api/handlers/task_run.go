package handlers

import (
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

func (h *TaskRunHandler) Run(c *gin.Context) {
	id := c.Param("id")
	task, err := h.runner.Run(c.Request.Context(), id)
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
