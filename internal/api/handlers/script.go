package handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"bench-hub/internal/model"
	"bench-hub/internal/service"
)

type ScriptHandler struct {
	scripts *service.ScriptService
}

type scriptCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Content     string `json:"content" binding:"required"`
}

type scriptUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

func NewScriptHandler(scripts *service.ScriptService) *ScriptHandler {
	return &ScriptHandler{scripts: scripts}
}

func (h *ScriptHandler) List(c *gin.Context) {
	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("page_size"), 20)
	if page < 1 || pageSize < 1 || pageSize > 100 {
		model.JSON(c, http.StatusBadRequest, model.Fail(1000, "invalid params"))
		return
	}

	offset := (page - 1) * pageSize
	scripts, err := h.scripts.List(c.Request.Context(), pageSize, offset)
	if err != nil {
		model.JSON(c, http.StatusInternalServerError, model.Fail(9000, "internal error"))
		return
	}

	model.JSON(c, http.StatusOK, model.OK(gin.H{
		"items": scripts,
		"page":  page,
		"size":  pageSize,
	}))
}

func (h *ScriptHandler) Get(c *gin.Context) {
	id := c.Param("id")
	script, err := h.scripts.Get(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrNotFound {
			model.JSON(c, http.StatusNotFound, model.Fail(1003, "not found"))
			return
		}
		model.JSON(c, http.StatusInternalServerError, model.Fail(9000, "internal error"))
		return
	}
	model.JSON(c, http.StatusOK, model.OK(script))
}

func (h *ScriptHandler) Create(c *gin.Context) {
	var req scriptCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		model.JSON(c, http.StatusBadRequest, model.Fail(1000, "invalid params"))
		return
	}

	script, err := h.scripts.Create(c.Request.Context(), req.Name, req.Description, req.Content)
	if err != nil {
		model.JSON(c, http.StatusInternalServerError, model.Fail(9000, "internal error"))
		return
	}

	model.JSON(c, http.StatusOK, model.OK(script))
}

func (h *ScriptHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req scriptUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		model.JSON(c, http.StatusBadRequest, model.Fail(1000, "invalid params"))
		return
	}

	script, err := h.scripts.Update(c.Request.Context(), id, req.Name, req.Description, req.Content)
	if err != nil {
		if err == service.ErrNotFound {
			model.JSON(c, http.StatusNotFound, model.Fail(1003, "not found"))
			return
		}
		model.JSON(c, http.StatusInternalServerError, model.Fail(9000, "internal error"))
		return
	}

	model.JSON(c, http.StatusOK, model.OK(script))
}

func (h *ScriptHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.scripts.Delete(c.Request.Context(), id); err != nil {
		if err == service.ErrNotFound {
			model.JSON(c, http.StatusNotFound, model.Fail(1003, "not found"))
			return
		}
		model.JSON(c, http.StatusInternalServerError, model.Fail(9000, "internal error"))
		return
	}

	model.JSON(c, http.StatusOK, model.OK(nil))
}

func (h *ScriptHandler) Import(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		model.JSON(c, http.StatusBadRequest, model.Fail(1000, "invalid params"))
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		model.JSON(c, http.StatusBadRequest, model.Fail(1000, "invalid params"))
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		model.JSON(c, http.StatusInternalServerError, model.Fail(9000, "internal error"))
		return
	}

	description := c.PostForm("description")
	script, err := h.scripts.Create(c.Request.Context(), name, description, string(data))
	if err != nil {
		model.JSON(c, http.StatusInternalServerError, model.Fail(9000, "internal error"))
		return
	}

	model.JSON(c, http.StatusOK, model.OK(script))
}
