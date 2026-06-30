package application

import "context"

type Template struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"-"`
}

type CreateServiceInput struct {
	Name        string
	Description string
	ProjectID   string
	TemplateID  string
}

type CreateServiceOutput struct {
	ServiceID     string
	RepositoryURL string
}

type ServiceOutput struct {
	ID          string
	Name        string
	Description string
	ProjectID   string
	TemplateID  string
	Repository  string
}

type UpdateServiceInput struct {
	ID          string
	Name        string
	Description string
}

type ServiceUsecase interface {
	CreateService(ctx context.Context, input CreateServiceInput) (CreateServiceOutput, error)
	GetServiceByID(ctx context.Context, id string) (ServiceOutput, error)
	ListProjectServices(ctx context.Context, projectID string) ([]ServiceOutput, error)
	UpdateService(ctx context.Context, input UpdateServiceInput) error
	DeleteService(ctx context.Context, id string) error

	ListTemplates(ctx context.Context) ([]Template, error)
}
