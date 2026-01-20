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
		token := ""
		header := c.GetHeader("Authorization")
		if strings.HasPrefix(header, "Bearer ") {
			token = strings.TrimPrefix(header, "Bearer ")
		}
		if token == "" {
			token = c.Query("token")
		}
		if token == "" {
			if cookie, err := c.Cookie("access_token"); err == nil {
				token = cookie
			}
		}
		if token == "" {
			model.JSON(c, http.StatusUnauthorized, model.Fail(1001, "unauthorized"))
			c.Abort()
			return
		}

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
