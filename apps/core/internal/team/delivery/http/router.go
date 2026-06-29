package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterTeamRoutes(r *gin.RouterGroup, handler *TeamHandler, authMiddleware gin.HandlerFunc) {
	teams := r.Group("/teams", authMiddleware)
	{
		teams.POST("", handler.Create)
		teams.GET("/:id", handler.GetByID)
		teams.GET("", handler.ListByOrg)
	}
}
