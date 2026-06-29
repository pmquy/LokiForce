package application

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
	CreateProject(input CreateProjectInput) (CreateProjectOutput, error)
	GetProjectByID(id string) (ProjectOutput, error)
	ListOrgProjects(orgID string) ([]ProjectOutput, error)
}
