package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"bench-hub/internal/model"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				model.JSON(c, http.StatusInternalServerError, model.Fail(9000, "internal error"))
				c.Abort()
			}
		}()
		c.Next()
	}
}
