package git

import (
	"context"
)

// GitClient defines the common interface for version control providers (GitHub, GitLab, etc.)
type GitClient interface {
	CreateRepository(ctx context.Context, name, description string, private bool) (string, error)
	DeleteRepository(ctx context.Context, name string) error
	GetFileSHA(ctx context.Context, repo, path string) (string, error)
	WriteFile(ctx context.Context, repo, path, content, message string) error
	DeleteFile(ctx context.Context, repo, path, message string) error
	PushFiles(ctx context.Context, repoURL string, files map[string][]byte, commitMessage string) error
	CreateRepositorySecret(ctx context.Context, repo, secretName, secretValue string) error
}
