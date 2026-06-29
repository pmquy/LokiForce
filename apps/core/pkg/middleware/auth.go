package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"lokiforce.com/apps/core/pkg/response"
)

type TokenValidator interface {
	ValidateToken(token string) (string, error)
}

func AuthMiddleware(validator TokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Fail(c, http.StatusUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Fail(c, http.StatusUnauthorized, "Authorization header must be Bearer <token>")
			c.Abort()
			return
		}

		token := parts[1]
		userID, err := validator.ValidateToken(token)
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, "Invalid token: "+err.Error())
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
