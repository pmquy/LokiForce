package application

import "context"

type CreateProjectInput struct {
	Name        string
	Description string
	OrgID       string
}

type CreateProjectOutput struct {
	ProjectID string
}

type ProjectOutput struct {
	ID          string
	Name        string
	Description string
	OrgID       string
}

type ProjectUsecase interface {
	CreateProject(ctx context.Context, input CreateProjectInput) (CreateProjectOutput, error)
	GetProjectByID(ctx context.Context, id string) (ProjectOutput, error)
	ListOrgProjects(ctx context.Context, orgID string) ([]ProjectOutput, error)
}
