package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"bench-hub/internal/model"
)

func Health(c *gin.Context) {
	model.JSON(c, http.StatusOK, model.OK(nil))
}

func Ping(c *gin.Context) {
	model.JSON(c, http.StatusOK, model.OK(gin.H{"ping": "pong"}))
}
