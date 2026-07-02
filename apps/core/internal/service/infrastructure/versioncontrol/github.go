package versioncontrol

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"text/template"

	"lokiforce.com/apps/core/internal/config"
	"lokiforce.com/apps/core/internal/service/application"
	gitclient "lokiforce.com/apps/core/pkg/git"
)

type GitHubVersionControl struct {
	client     gitclient.GitClient
	owner      string
	gitopsRepo string
	token      string
}

func NewGitHubVersionControl(cfg *config.Config, client gitclient.GitClient) application.VersionControl {
	return &GitHubVersionControl{
		client:     client,
		owner:      cfg.GitHub.Owner,
		gitopsRepo: cfg.GitHub.GitOpsRepo,
		token:      cfg.GitHub.Token,
	}
}

func (g *GitHubVersionControl) CreateRepository(ctx context.Context, config application.RepositoryConfig) (string, error) {
	slog.Info("Creating remote Git repository via client", "owner", g.owner, "repo", config.Name)
	cloneURL, err := g.client.CreateRepository(ctx, config.Name, config.Description, config.Private)
	if err != nil {
		return "", err
	}

	slog.Info("Injecting GITOPS_GIT_TOKEN secret into remote repository", "repo", config.Name)
	if err := g.client.CreateRepositorySecret(ctx, config.Name, "GITOPS_GIT_TOKEN", g.token); err != nil {
		slog.Error("Failed to inject repository secret GITOPS_GIT_TOKEN", "repo", config.Name, "error", err)
	}

	return cloneURL, nil
}

func (g *GitHubVersionControl) PushFiles(ctx context.Context, repoURL string, templateDir string, serviceName string) error {
	slog.Info("Rendering template files for remote repository", "dir", templateDir)

	filesToCommit := make(map[string][]byte)

	var renderDir func(src string, dst string) error
	renderDir = func(src string, dst string) error {
		entries, err := os.ReadDir(src)
		if err != nil {
			return err
		}

		data := map[string]string{
			"ServiceName": serviceName,
			"Owner":       g.owner,
			"GitOpsRepo":  g.gitopsRepo,
		}

		for _, entry := range entries {
			srcPath := filepath.Join(src, entry.Name())
			dstPath := filepath.Join(dst, entry.Name())

			if entry.IsDir() {
				if entry.Name() == "deploy" {
					slog.Info("Skipping deploy folder in template copying", "path", srcPath)
					continue
				}
				if err := renderDir(srcPath, dstPath); err != nil {
					return err
				}
			} else {
				fileBytes, err := os.ReadFile(srcPath)
				if err != nil {
					return err
				}

				tmpl, err := template.New(entry.Name()).Parse(string(fileBytes))
				if err != nil {
					return fmt.Errorf("failed to parse template file %s: %w", entry.Name(), err)
				}

				var buf bytes.Buffer
				if err := tmpl.Execute(&buf, data); err != nil {
					return fmt.Errorf("failed to execute template file %s: %w", entry.Name(), err)
				}

				filesToCommit[dstPath] = buf.Bytes()
			}
		}
		return nil
	}

	if err := renderDir(templateDir, ""); err != nil {
		return err
	}

	commitMessage := "chore(scaffolding): initialize service template from LokiForce Portal"
	return g.client.PushFiles(ctx, repoURL, filesToCommit, commitMessage)
}

func (g *GitHubVersionControl) DeleteRepository(ctx context.Context, name string) error {
	slog.Info("Deleting remote Git repository via client", "owner", g.owner, "repo", name)
	return g.client.DeleteRepository(ctx, name)
}
