package application

import (
	"context"

	"github.com/google/uuid"
	"lokiforce.com/apps/core/internal/project/domain"
)

type projectUsecaseImpl struct {
	repository domain.ProjectRepository
}

func NewProjectUsecase(repo domain.ProjectRepository) ProjectUsecase {
	return &projectUsecaseImpl{repository: repo}
}

func (u *projectUsecaseImpl) CreateProject(ctx context.Context, input CreateProjectInput) (CreateProjectOutput, error) {
	id := uuid.NewString()
	project, err := domain.NewProject(id, input.Name, input.Description, input.OrgID)
	if err != nil {
		return CreateProjectOutput{}, err
	}

	if err := u.repository.Create(ctx, project); err != nil {
		return CreateProjectOutput{}, err
	}

	return CreateProjectOutput{ProjectID: project.ID}, nil
}

func (u *projectUsecaseImpl) GetProjectByID(ctx context.Context, id string) (ProjectOutput, error) {
	project, err := u.repository.GetByID(ctx, id)
	if err != nil {
		return ProjectOutput{}, err
	}

	return ProjectOutput{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		OrgID:       project.OrgID,
	}, nil
}

func (u *projectUsecaseImpl) ListOrgProjects(ctx context.Context, orgID string) ([]ProjectOutput, error) {
	projects, err := u.repository.ListByOrg(ctx, orgID)
	if err != nil {
		return nil, err
	}

	outputs := make([]ProjectOutput, len(projects))
	for i, project := range projects {
		outputs[i] = ProjectOutput{
			ID:          project.ID,
			Name:        project.Name,
			Description: project.Description,
			OrgID:       project.OrgID,
		}
	}
	return outputs, nil
}
