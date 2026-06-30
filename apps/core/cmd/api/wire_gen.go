//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject

package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lokiforce.com/apps/core/internal/config"
	application2 "lokiforce.com/apps/core/internal/organization/application"
	http2 "lokiforce.com/apps/core/internal/organization/delivery/http"
	repository2 "lokiforce.com/apps/core/internal/organization/infrastructure/repository"
	application3 "lokiforce.com/apps/core/internal/project/application"
	http3 "lokiforce.com/apps/core/internal/project/delivery/http"
	repository3 "lokiforce.com/apps/core/internal/project/infrastructure/repository"
	application5 "lokiforce.com/apps/core/internal/service/application"
	http5 "lokiforce.com/apps/core/internal/service/delivery/http"
	repository5 "lokiforce.com/apps/core/internal/service/infrastructure/repository"
	"lokiforce.com/apps/core/internal/service/infrastructure/versioncontrol"
	application4 "lokiforce.com/apps/core/internal/team/application"
	http4 "lokiforce.com/apps/core/internal/team/delivery/http"
	repository4 "lokiforce.com/apps/core/internal/team/infrastructure/repository"
	"lokiforce.com/apps/core/internal/user/application"
	"lokiforce.com/apps/core/internal/user/delivery/http"
	"lokiforce.com/apps/core/internal/user/infrastructure/jwt"
	"lokiforce.com/apps/core/internal/user/infrastructure/repository"
)

func InitializeApp(cfg *config.Config) (*Handlers, error) {
	db, err := ProvideDB(cfg)
	if err != nil {
		return nil, err
	}
	postgresUserRepository := repository.NewPostgresUserRepository(db)
	tokenService := ProvideTokenService(cfg)
	userUsecase := application.NewUserUsecase(postgresUserRepository, tokenService)
	userHandler := http.NewUserHandler(userUsecase)
	postgresOrgRepository := repository2.NewPostgresOrgRepository(db)
	orgUsecase := application2.NewOrgUsecase(postgresOrgRepository)
	orgHandler := http2.NewOrgHandler(orgUsecase)
	postgresProjectRepository := repository3.NewPostgresProjectRepository(db)
	projectUsecase := application3.NewProjectUsecase(postgresProjectRepository)
	projectHandler := http3.NewProjectHandler(projectUsecase)
	postgresTeamRepository := repository4.NewPostgresTeamRepository(db)
	teamUsecase := application4.NewTeamUsecase(postgresTeamRepository)
	teamHandler := http4.NewTeamHandler(teamUsecase)
	postgresServiceRepository := repository5.NewPostgresServiceRepository(db)
	versionControl := versioncontrol.NewGitHubVersionControl(cfg)
	serviceUsecase := application5.NewServiceUsecase(postgresServiceRepository, versionControl)
	serviceHandler := http5.NewServiceHandler(serviceUsecase)
	handlers := NewHandlers(userHandler, orgHandler, projectHandler, teamHandler, serviceHandler)
	return handlers, nil
}

func ProvideDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.GetDatabaseURL()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := repository.Migrate(db); err != nil {
		return nil, err
	}
	if err := repository2.Migrate(db); err != nil {
		return nil, err
	}
	if err := repository3.Migrate(db); err != nil {
		return nil, err
	}
	if err := repository4.Migrate(db); err != nil {
		return nil, err
	}
	if err := repository5.Migrate(db); err != nil {
		return nil, err
	}
	return db, nil
}

func ProvideTokenService(cfg *config.Config) application.TokenService {
	return jwt.NewJWTService(cfg.JWT.Secret)
}
