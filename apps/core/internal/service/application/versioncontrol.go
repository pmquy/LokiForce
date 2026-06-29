package application

import "context"

type RepositoryConfig struct {
	Name        string
	Description string
	Private     bool
}

type VersionControl interface {
	CreateRepository(ctx context.Context, config RepositoryConfig) (string, error)
	PushFiles(ctx context.Context, repoURL string, localDir string) error
}
