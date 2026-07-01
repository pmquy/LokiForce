package application

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"lokiforce.com/apps/core/internal/service/domain"
)

var AvailableTemplates = []Template{
	{
		ID:          "golang",
		Name:        "Golang Gin API Template",
		Description: "Standard Golang microservice template built with Gin Gonic.",
		Path:        "./templates/golang",
	},
	{
		ID:          "nodejs",
		Name:        "NodeJS Express API Template",
		Description: "Standard Node.js microservice template built with Express.",
		Path:        "./templates/nodejs",
	},
}

type serviceUsecaseImpl struct {
	repository        domain.ServiceRepository
	versionControl    VersionControl
	deploymentControl DeploymentControl
}

func NewServiceUsecase(
	repo domain.ServiceRepository,
	vc VersionControl,
	dc DeploymentControl,
) ServiceUsecase {
	return &serviceUsecaseImpl{
		repository:        repo,
		versionControl:    vc,
		deploymentControl: dc,
	}
}

func (u *serviceUsecaseImpl) CreateService(ctx context.Context, input CreateServiceInput) (CreateServiceOutput, error) {

	var matchTemplate *Template
	for _, t := range AvailableTemplates {
		if t.ID == input.TemplateID {
			matchTemplate = &t
			break
		}
	}
	if matchTemplate == nil {
		return CreateServiceOutput{}, fmt.Errorf("invalid template ID: %s", input.TemplateID)
	}

	id := uuid.NewString()
	svc, err := domain.NewService(id, input.Name, input.Description, input.ProjectID, input.TemplateID)
	if err != nil {
		return CreateServiceOutput{}, err
	}

	if err := u.repository.Create(ctx, svc); err != nil {
		return CreateServiceOutput{}, err
	}

	uniqueRepoName := fmt.Sprintf("%s-%s", svc.Name, svc.ID[:8])
	slog.Info("Creating Git repository", "name", uniqueRepoName)
	repoConfig := RepositoryConfig{
		Name:        uniqueRepoName,
		Description: svc.Description,
		Private:     false,
	}
	repoURL, err := u.versionControl.CreateRepository(ctx, repoConfig)
	if err != nil {
		return CreateServiceOutput{}, fmt.Errorf("failed to create Git repository: %w", err)
	}

	slog.Info("Scaffolding and pushing files to Git repository in-memory", "url", repoURL)
	if err := u.versionControl.PushFiles(ctx, repoURL, matchTemplate.Path, uniqueRepoName); err != nil {
		return CreateServiceOutput{}, fmt.Errorf("failed to push template code: %w", err)
	}

	slog.Info("Registering GitOps Deployment (Argo CD)", "serviceName", uniqueRepoName)
	deployConfig := DeploymentConfig{
		ServiceName:   uniqueRepoName,
		RepositoryURL: repoURL,
		Namespace:     "production",
		Environment:   input.ProjectID,
	}
	_, err = u.deploymentControl.RegisterDeployment(ctx, deployConfig)
	if err != nil {
		return CreateServiceOutput{}, fmt.Errorf("failed to register deployment: %w", err)
	}

	svc.Repository = repoURL
	if err := u.repository.Update(ctx, svc); err != nil {
		return CreateServiceOutput{}, fmt.Errorf("failed to update service repo url: %w", err)
	}

	return CreateServiceOutput{
		ServiceID:     svc.ID,
		RepositoryURL: repoURL,
	}, nil
}

func (u *serviceUsecaseImpl) GetServiceByID(ctx context.Context, id string) (ServiceOutput, error) {
	svc, err := u.repository.GetByID(ctx, id)
	if err != nil {
		return ServiceOutput{}, err
	}

	return ServiceOutput{
		ID:          svc.ID,
		Name:        svc.Name,
		Description: svc.Description,
		ProjectID:   svc.ProjectID,
		TemplateID:  svc.TemplateID,
		Repository:  svc.Repository,
	}, nil
}

func (u *serviceUsecaseImpl) ListProjectServices(ctx context.Context, projectID string) ([]ServiceOutput, error) {
	services, err := u.repository.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	outputs := make([]ServiceOutput, len(services))
	for i, svc := range services {
		outputs[i] = ServiceOutput{
			ID:          svc.ID,
			Name:        svc.Name,
			Description: svc.Description,
			ProjectID:   svc.ProjectID,
			TemplateID:  svc.TemplateID,
			Repository:  svc.Repository,
		}
	}
	return outputs, nil
}

func (u *serviceUsecaseImpl) UpdateService(ctx context.Context, input UpdateServiceInput) error {
	svc, err := u.repository.GetByID(ctx, input.ID)
	if err != nil {
		return err
	}

	svc.Name = input.Name
	svc.Description = input.Description

	return u.repository.Update(ctx, svc)
}

func (u *serviceUsecaseImpl) DeleteService(ctx context.Context, id string) error {
	return u.repository.Delete(ctx, id)
}

func (u *serviceUsecaseImpl) ListTemplates(ctx context.Context) ([]Template, error) {
	return AvailableTemplates, nil
}
