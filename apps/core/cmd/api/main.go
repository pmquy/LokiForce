package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"lokiforce.com/apps/core/internal/config"
	orgHttp "lokiforce.com/apps/core/internal/organization/delivery/http"
	projectHttp "lokiforce.com/apps/core/internal/project/delivery/http"
	teamHttp "lokiforce.com/apps/core/internal/team/delivery/http"
	userHttp "lokiforce.com/apps/core/internal/user/delivery/http"
	"lokiforce.com/apps/core/pkg/middleware"
)

type Handlers struct {
	UserHandler *userHttp.UserHandler
	OrgHandler  *orgHttp.OrgHandler
	ProjHandler *projectHttp.ProjectHandler
	TeamHandler *teamHttp.TeamHandler
}

func NewHandlers(
	user *userHttp.UserHandler,
	org *orgHttp.OrgHandler,
	proj *projectHttp.ProjectHandler,
	team *teamHttp.TeamHandler,
) *Handlers {
	return &Handlers{
		UserHandler: user,
		OrgHandler:  org,
		ProjHandler: proj,
		TeamHandler: team,
	}
}

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize app using wire injection
	handlers, err := InitializeApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	// Set up router
	r := gin.Default()

	apiV1 := r.Group("/api/v1")
	tokenService := ProvideTokenService(cfg)
	authMiddleware := middleware.AuthMiddleware(tokenService)

	userHttp.RegisterUserRoutes(apiV1, handlers.UserHandler, tokenService)
	orgHttp.RegisterOrgRoutes(apiV1, handlers.OrgHandler, authMiddleware)
	teamHttp.RegisterTeamRoutes(apiV1, handlers.TeamHandler, authMiddleware)
	projectHttp.RegisterProjectRoutes(apiV1, handlers.ProjHandler, authMiddleware)

	portStr := fmt.Sprintf("%d", cfg.Server.Port)
	log.Printf("Starting server on :%s...\n", portStr)
	if err := r.Run(":" + portStr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
