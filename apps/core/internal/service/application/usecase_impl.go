package application

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

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
	repository     domain.ServiceRepository
	versionControl VersionControl
}

func NewServiceUsecase(repo domain.ServiceRepository, vc VersionControl) ServiceUsecase {
	return &serviceUsecaseImpl{
		repository:     repo,
		versionControl: vc,
	}
}

func (u *serviceUsecaseImpl) CreateService(ctx context.Context, input CreateServiceInput) (CreateServiceOutput, error) {
	// Validate template ID
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

	// 1. Save metadata to DB
	if err := u.repository.Create(ctx, svc); err != nil {
		return CreateServiceOutput{}, err
	}

	// 2. Generate scaffold files locally in a temporary directory
	tempDir := filepath.Join("./scratch", "temp_"+svc.ID)
	outputPath := filepath.Join(tempDir, svc.Name)
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return CreateServiceOutput{}, fmt.Errorf("failed to create temporary folder: %w", err)
	}

	// Clean up local temp folder after operations
	defer os.RemoveAll(tempDir)

	if err := copyDir(matchTemplate.Path, outputPath); err != nil {
		return CreateServiceOutput{}, fmt.Errorf("failed to copy template files: %w", err)
	}

		// 3. Create remote repository on Git provider
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

	// 4. Push generated files to the remote repository
	slog.Info("Pushing files to Git repository", "url", repoURL)
	if err := u.versionControl.PushFiles(ctx, repoURL, outputPath); err != nil {
		return CreateServiceOutput{}, fmt.Errorf("failed to push template code: %w", err)
	}

	// 5. Update service database record with the repository URL
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

// copyDir recursively copies a directory tree.
func copyDir(src string, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return err
			}
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			data, err := os.ReadFile(srcPath)
			if err != nil {
				return err
			}
			if err := os.WriteFile(dstPath, data, 0644); err != nil {
				return err
			}
		}
	}
	return nil
}
