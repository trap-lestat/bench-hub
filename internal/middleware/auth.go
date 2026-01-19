package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"bench-hub/internal/model"
	"bench-hub/internal/service"
)

func Auth(auth *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			model.JSON(c, http.StatusUnauthorized, model.Fail(1001, "unauthorized"))
			c.Abort()
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		userID, err := auth.ValidateAccess(token)
		if err != nil {
			model.JSON(c, http.StatusUnauthorized, model.Fail(1001, "unauthorized"))
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
