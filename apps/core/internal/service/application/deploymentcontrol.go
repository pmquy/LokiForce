package application

import "context"

type DeploymentConfig struct {
	ServiceName   string
	RepositoryURL string
	Namespace     string
	Environment   string
	TemplateID    string
}

type DeploymentControl interface {
	RegisterDeployment(ctx context.Context, config DeploymentConfig) (string, error)
	DeleteDeployment(ctx context.Context, serviceName string) error
}
