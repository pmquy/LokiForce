//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lokiforce.com/apps/core/internal/config"

	// Organization
	orgApp "lokiforce.com/apps/core/internal/organization/application"
	orgHttp "lokiforce.com/apps/core/internal/organization/delivery/http"
	orgDomain "lokiforce.com/apps/core/internal/organization/domain"
	orgRepo "lokiforce.com/apps/core/internal/organization/infrastructure/repository"

	// Project
	projectApp "lokiforce.com/apps/core/internal/project/application"
	projectHttp "lokiforce.com/apps/core/internal/project/delivery/http"
	projectDomain "lokiforce.com/apps/core/internal/project/domain"
	projectRepo "lokiforce.com/apps/core/internal/project/infrastructure/repository"

	// Service
	serviceApp "lokiforce.com/apps/core/internal/service/application"
	serviceHttp "lokiforce.com/apps/core/internal/service/delivery/http"
	serviceDomain "lokiforce.com/apps/core/internal/service/domain"
	serviceRepo "lokiforce.com/apps/core/internal/service/infrastructure/repository"
	serviceVc "lokiforce.com/apps/core/internal/service/infrastructure/versioncontrol"

	// Team
	teamApp "lokiforce.com/apps/core/internal/team/application"
	teamHttp "lokiforce.com/apps/core/internal/team/delivery/http"
	teamDomain "lokiforce.com/apps/core/internal/team/domain"
	teamRepo "lokiforce.com/apps/core/internal/team/infrastructure/repository"

	// User
	"lokiforce.com/apps/core/internal/user/application"
	userHttp "lokiforce.com/apps/core/internal/user/delivery/http"
	"lokiforce.com/apps/core/internal/user/domain"
	"lokiforce.com/apps/core/internal/user/infrastructure/jwt"
	"lokiforce.com/apps/core/internal/user/infrastructure/repository"
)

func ProvideDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.GetDatabaseURL()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := repository.Migrate(db); err != nil {
		return nil, err
	}
	if err := orgRepo.Migrate(db); err != nil {
		return nil, err
	}
	if err := projectRepo.Migrate(db); err != nil {
		return nil, err
	}
	if err := teamRepo.Migrate(db); err != nil {
		return nil, err
	}
	if err := serviceRepo.Migrate(db); err != nil {
		return nil, err
	}
	return db, nil
}

func ProvideTokenService(cfg *config.Config) application.TokenService {
	return jwt.NewJWTService(cfg.JWT.Secret)
}

func InitializeApp(cfg *config.Config) (*Handlers, error) {
	wire.Build(
		ProvideDB,
		ProvideTokenService,

		// User
		repository.NewPostgresUserRepository,
		wire.Bind(new(domain.UserRepository), new(*repository.PostgresUserRepository)),
		application.NewUserUsecase,
		userHttp.NewUserHandler,

		// Organization
		orgRepo.NewPostgresOrgRepository,
		wire.Bind(new(orgDomain.OrganizationRepository), new(*orgRepo.PostgresOrgRepository)),
		orgApp.NewOrgUsecase,
		orgHttp.NewOrgHandler,

		// Project
		projectRepo.NewPostgresProjectRepository,
		wire.Bind(new(projectDomain.ProjectRepository), new(*projectRepo.PostgresProjectRepository)),
		projectApp.NewProjectUsecase,
		projectHttp.NewProjectHandler,

		// Team
		teamRepo.NewPostgresTeamRepository,
		wire.Bind(new(teamDomain.TeamRepository), new(*teamRepo.PostgresTeamRepository)),
		teamApp.NewTeamUsecase,
		teamHttp.NewTeamHandler,

		// Service
		serviceRepo.NewPostgresServiceRepository,
		wire.Bind(new(serviceDomain.ServiceRepository), new(*serviceRepo.PostgresServiceRepository)),
		serviceVc.NewGitHubVersionControl,
		serviceApp.NewServiceUsecase,
		serviceHttp.NewServiceHandler,

		NewHandlers,
	)
	return nil, nil
}
