package application

import "context"

type RepositoryConfig struct {
	Name        string
	Description string
	Private     bool
}

type VersionControl interface {
	CreateRepository(ctx context.Context, config RepositoryConfig) (string, error)
	PushFiles(ctx context.Context, repoURL string, templateDir string, serviceName string) error
	DeleteRepository(ctx context.Context, name string) error
}
