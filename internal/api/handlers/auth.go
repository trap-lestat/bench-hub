package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"bench-hub/internal/model"
	"bench-hub/internal/service"
)

type AuthHandler struct {
	auth *service.AuthService
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		model.JSON(c, http.StatusBadRequest, model.Fail(1000, "invalid params"))
		return
	}

	access, refresh, user, err := h.auth.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			model.JSON(c, http.StatusUnauthorized, model.Fail(1001, "invalid credentials"))
			return
		}
		model.JSON(c, http.StatusInternalServerError, model.Fail(9000, "internal error"))
		return
	}

	model.JSON(c, http.StatusOK, model.OK(gin.H{
		"access_token":  access,
		"refresh_token": refresh,
		"user":          user,
	}))
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		model.JSON(c, http.StatusBadRequest, model.Fail(1000, "invalid params"))
		return
	}

	access, refresh, err := h.auth.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		model.JSON(c, http.StatusUnauthorized, model.Fail(1001, "unauthorized"))
		return
	}

	model.JSON(c, http.StatusOK, model.OK(gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	}))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	model.JSON(c, http.StatusOK, model.OK(nil))
}
