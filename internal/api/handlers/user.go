package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"bench-hub/internal/model"
	"bench-hub/internal/service"
)

type UserHandler struct {
	users *service.UserService
}

type userCreateRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type userUpdateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUserHandler(users *service.UserService) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) List(c *gin.Context) {
	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("page_size"), 20)
	if page < 1 || pageSize < 1 || pageSize > 100 {
		model.JSON(c, http.StatusBadRequest, model.Fail(1000, "invalid params"))
		return
	}

	offset := (page - 1) * pageSize
	users, err := h.users.List(c.Request.Context(), pageSize, offset)
	if err != nil {
		model.JSON(c, http.StatusInternalServerError, model.Fail(9000, "internal error"))
		return
	}

	model.JSON(c, http.StatusOK, model.OK(gin.H{
		"items": users,
		"page":  page,
		"size":  pageSize,
	}))
}

func (h *UserHandler) Get(c *gin.Context) {
	id := c.Param("id")
	user, err := h.users.Get(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrNotFound {
			model.JSON(c, http.StatusNotFound, model.Fail(1003, "not found"))
			return
		}
		model.JSON(c, http.StatusInternalServerError, model.Fail(9000, "internal error"))
		return
	}
	model.JSON(c, http.StatusOK, model.OK(user))
}

func (h *UserHandler) Create(c *gin.Context) {
	var req userCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		model.JSON(c, http.StatusBadRequest, model.Fail(1000, "invalid params"))
		return
	}

	user, err := h.users.Create(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		model.JSON(c, http.StatusInternalServerError, model.Fail(9000, "internal error"))
		return
	}

	model.JSON(c, http.StatusOK, model.OK(user))
}

func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req userUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		model.JSON(c, http.StatusBadRequest, model.Fail(1000, "invalid params"))
		return
	}

	user, err := h.users.Update(c.Request.Context(), id, req.Username, req.Password)
	if err != nil {
		if err == service.ErrNotFound {
			model.JSON(c, http.StatusNotFound, model.Fail(1003, "not found"))
			return
		}
		model.JSON(c, http.StatusInternalServerError, model.Fail(9000, "internal error"))
		return
	}

	model.JSON(c, http.StatusOK, model.OK(user))
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.users.Delete(c.Request.Context(), id); err != nil {
		if err == service.ErrNotFound {
			model.JSON(c, http.StatusNotFound, model.Fail(1003, "not found"))
			return
		}
		model.JSON(c, http.StatusInternalServerError, model.Fail(9000, "internal error"))
		return
	}

	model.JSON(c, http.StatusOK, model.OK(nil))
}
