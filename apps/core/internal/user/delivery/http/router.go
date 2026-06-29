package http

import (
	"github.com/gin-gonic/gin"
	"lokiforce.com/apps/core/internal/user/application"
)

func RegisterUserRoutes(r *gin.RouterGroup, handler *UserHandler, authMiddleware gin.HandlerFunc, tokenService application.TokenService) {
	users := r.Group("/users")
	{
		users.POST("/register", handler.register)
		users.POST("/login", handler.login)
		users.GET("/profile", authMiddleware, handler.getProfile)
	}
}
