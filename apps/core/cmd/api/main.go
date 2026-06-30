package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"lokiforce.com/apps/core/internal/config"
	orgHttp "lokiforce.com/apps/core/internal/organization/delivery/http"
	projectHttp "lokiforce.com/apps/core/internal/project/delivery/http"
	serviceHttp "lokiforce.com/apps/core/internal/service/delivery/http"
	teamHttp "lokiforce.com/apps/core/internal/team/delivery/http"
	userHttp "lokiforce.com/apps/core/internal/user/delivery/http"
	"lokiforce.com/apps/core/pkg/middleware"
)

type Handlers struct {
	UserHandler *userHttp.UserHandler
	OrgHandler  *orgHttp.OrgHandler
	ProjHandler *projectHttp.ProjectHandler
	TeamHandler *teamHttp.TeamHandler
	SvcHandler  *serviceHttp.ServiceHandler
}

func NewHandlers(
	user *userHttp.UserHandler,
	org *orgHttp.OrgHandler,
	proj *projectHttp.ProjectHandler,
	team *teamHttp.TeamHandler,
	svc *serviceHttp.ServiceHandler,
) *Handlers {
	return &Handlers{
		UserHandler: user,
		OrgHandler:  org,
		ProjHandler: proj,
		TeamHandler: team,
		SvcHandler:  svc,
	}
}

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	handlers, err := InitializeApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	r := gin.Default()
	r.Use(middleware.CORSMiddleware(cfg.Server.AllowedOrigins))

	apiV1 := r.Group("/api/v1")
	tokenService := ProvideTokenService(cfg)
	authMiddleware := middleware.AuthMiddleware(tokenService)

	userHttp.RegisterUserRoutes(apiV1, handlers.UserHandler, authMiddleware, tokenService)
	orgHttp.RegisterOrgRoutes(apiV1, handlers.OrgHandler, authMiddleware)
	projectHttp.RegisterProjectRoutes(apiV1, handlers.ProjHandler, authMiddleware)
	teamHttp.RegisterTeamRoutes(apiV1, handlers.TeamHandler, authMiddleware)
	serviceHttp.RegisterServiceRoutes(apiV1, handlers.SvcHandler, authMiddleware)

	portStr := fmt.Sprintf("%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    ":" + portStr,
		Handler: r,
	}

	go func() {
		slog.Info("Starting server", "port", portStr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exited cleanly")
}
