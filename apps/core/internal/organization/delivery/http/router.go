package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterOrgRoutes(r *gin.RouterGroup, handler *OrgHandler, authMiddleware gin.HandlerFunc) {
	orgs := r.Group("/organizations", authMiddleware)
	{
		orgs.POST("", handler.Create)
		orgs.GET("/:id", handler.GetByID)
		orgs.GET("", handler.ListMyOrgs)
	}
}
