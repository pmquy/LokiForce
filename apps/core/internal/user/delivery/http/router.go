package http

import (
	"github.com/gin-gonic/gin"
	"lokiforce.com/apps/core/internal/user/application"
	"lokiforce.com/apps/core/pkg/middleware"
)

func RegisterUserRoutes(r *gin.RouterGroup, handler *UserHandler, tokenService application.TokenService) {
	users := r.Group("/users")
	{
		users.POST("/register", handler.register)
		users.POST("/login", handler.login)
		users.GET("/profile", middleware.AuthMiddleware(tokenService), handler.getProfile)
	}
}
