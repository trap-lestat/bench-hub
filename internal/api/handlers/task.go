package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"bench-hub/internal/model"
	"bench-hub/internal/service"
)

type TaskHandler struct {
	tasks *service.TaskService
}

type taskCreateRequest struct {
	Name            string `json:"name" binding:"required"`
	ScriptID        string `json:"script_id" binding:"required"`
	UsersCount      int    `json:"users_count" binding:"required"`
	SpawnRate       int    `json:"spawn_rate" binding:"required"`
	DurationSeconds int    `json:"duration_seconds" binding:"required"`
}

func NewTaskHandler(tasks *service.TaskService) *TaskHandler {
	return &TaskHandler{tasks: tasks}
}

func (h *TaskHandler) List(c *gin.Context) {
	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("page_size"), 20)
	if page < 1 || pageSize < 1 || pageSize > 100 {
		model.JSON(c, http.StatusBadRequest, model.Fail(1000, "invalid params"))
		return
	}

	offset := (page - 1) * pageSize
	tasks, err := h.tasks.List(c.Request.Context(), pageSize, offset)
	if err != nil {
		model.JSON(c, http.StatusInternalServerError, model.Fail(9000, "internal error"))
		return
	}

	model.JSON(c, http.StatusOK, model.OK(gin.H{
		"items": tasks,
		"page":  page,
		"size":  pageSize,
	}))
}

func (h *TaskHandler) Get(c *gin.Context) {
	id := c.Param("id")
	task, err := h.tasks.Get(c.Request.Context(), id)
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

func (h *TaskHandler) Create(c *gin.Context) {
	var req taskCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		model.JSON(c, http.StatusBadRequest, model.Fail(1000, "invalid params"))
		return
	}

	if req.UsersCount <= 0 || req.SpawnRate <= 0 || req.DurationSeconds <= 0 {
		model.JSON(c, http.StatusBadRequest, model.Fail(1000, "invalid params"))
		return
	}

	task, err := h.tasks.Create(c.Request.Context(), req.Name, req.ScriptID, req.UsersCount, req.SpawnRate, req.DurationSeconds)
	if err != nil {
		model.JSON(c, http.StatusInternalServerError, model.Fail(9000, "internal error"))
		return
	}

	model.JSON(c, http.StatusOK, model.OK(task))
}

func (h *TaskHandler) Stop(c *gin.Context) {
	id := c.Param("id")
	task, err := h.tasks.Stop(c.Request.Context(), id)
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
