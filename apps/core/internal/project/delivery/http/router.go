package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterProjectRoutes(r *gin.RouterGroup, handler *ProjectHandler, authMiddleware gin.HandlerFunc) {
	projects := r.Group("/projects", authMiddleware)
	{
		projects.POST("", handler.Create)
		projects.GET("/:id", handler.GetByID)
		projects.GET("", handler.ListByOrg)
	}
}
