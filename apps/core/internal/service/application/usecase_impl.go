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

	repoName := getExternalName(svc)
	slog.Info("Creating Git repository", "name", repoName)
	repoConfig := RepositoryConfig{
		Name:        repoName,
		Description: svc.Description,
		Private:     false,
	}
	repoURL, err := u.versionControl.CreateRepository(ctx, repoConfig)
	if err != nil {
		return CreateServiceOutput{}, fmt.Errorf("failed to create Git repository: %w", err)
	}

	slog.Info("Resolved actual Git repository name", "actualName", repoName)

	slog.Info("Scaffolding and pushing files to Git repository in-memory", "url", repoURL)
	if err := u.versionControl.PushFiles(ctx, repoURL, matchTemplate.Path, repoName); err != nil {
		return CreateServiceOutput{}, fmt.Errorf("failed to push template code: %w", err)
	}

	slog.Info("Registering GitOps Deployment (Argo CD)", "serviceName", repoName)
	deployConfig := DeploymentConfig{
		ServiceName:   repoName,
		RepositoryURL: repoURL,
		Namespace:     getExternalNamespace(svc),
		Environment:   "production",
		TemplateID:    input.TemplateID,
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
	svc, err := u.repository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	repoName := getExternalName(svc)
	if repoName != "" {
		if err := u.deploymentControl.DeleteDeployment(ctx, repoName); err != nil {
			slog.Error("Failed to delete Argo CD deployment during service deletion", "repo", repoName, "error", err)
		}

		if err := u.versionControl.DeleteRepository(ctx, repoName); err != nil {
			slog.Error("Failed to delete GitHub repository during service deletion", "repo", repoName, "error", err)
		}
	}

	return u.repository.Delete(ctx, id)
}

func (u *serviceUsecaseImpl) ListTemplates(ctx context.Context) ([]Template, error) {
	return AvailableTemplates, nil
}

func (u *serviceUsecaseImpl) CreateAccessPolicy(ctx context.Context, clientID, targetID, targetPort, projectID string) (string, error) {
	policy, err := domain.NewAccessPolicy(uuid.NewString(), clientID, targetID, targetPort, projectID)
	if err != nil {
		return "", err
	}

	if err := u.repository.CreateAccessPolicy(ctx, policy); err != nil {
		return "", err
	}

	return policy.ID, nil
}

func (u *serviceUsecaseImpl) DeleteAccessPolicy(ctx context.Context, policyID string) error {
	return u.repository.DeleteAccessPolicy(ctx, policyID)
}

func getExternalName(svc *domain.Service) string {
	if svc == nil {
		return ""
	}
	return svc.ID
}

func getExternalNamespace(svc *domain.Service) string {
	if svc == nil {
		return ""
	}
	return svc.ProjectID
}
