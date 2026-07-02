package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterServiceRoutes(r *gin.RouterGroup, handler *ServiceHandler, authMiddleware gin.HandlerFunc) {
	services := r.Group("/services", authMiddleware)
	{
		services.POST("", handler.Create)
		services.GET("/:id", handler.GetByID)
		services.GET("", handler.ListByProject)
		services.PUT("/:id", handler.Update)
		services.DELETE("/:id", handler.Delete)
		services.GET("/templates", handler.ListTemplates)
		services.POST("/policies", handler.CreateAccessPolicy)
		services.DELETE("/policies/:policyId", handler.DeleteAccessPolicy)
		services.GET("/:id/policies", handler.ListAccessPolicies)
	}
}
