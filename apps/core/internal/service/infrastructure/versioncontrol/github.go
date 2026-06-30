package versioncontrol

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"lokiforce.com/apps/core/internal/config"
	"lokiforce.com/apps/core/internal/service/application"
)

type GitHubVersionControl struct {
	token  string
	owner  string
	prefix string
}

func NewGitHubVersionControl(cfg *config.Config) application.VersionControl {
	return &GitHubVersionControl{
		token:  cfg.GitHub.Token,
		owner:  cfg.GitHub.Owner,
		prefix: cfg.GitHub.Prefix,
	}
}

type createRepoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
}

type createRepoResponse struct {
	CloneURL string `json:"clone_url"`
}

func (g *GitHubVersionControl) CreateRepository(ctx context.Context, config application.RepositoryConfig) (string, error) {
	repoName := g.prefix + config.Name

	reqBody := createRepoRequest{
		Name:        repoName,
		Description: config.Description,
		Private:     config.Private,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	apiURL := "https://api.github.com/user/repos"
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("github api failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var respData createRepoResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", err
	}

	return respData.CloneURL, nil
}

func (g *GitHubVersionControl) PushFiles(ctx context.Context, repoURL string, templateDir string, serviceName string) error {
	fs := memfs.New()
	store := memory.NewStorage()

	r, err := git.Init(store, fs)
	if err != nil {
		return fmt.Errorf("failed to init in-memory git repository: %w", err)
	}

	if err := copyToMemFS(templateDir, fs, "", serviceName); err != nil {
		return fmt.Errorf("failed to load template files into memory: %w", err)
	}

	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get git worktree: %w", err)
	}

	_, err = w.Add(".")
	if err != nil {
		return fmt.Errorf("failed to add files to staging: %w", err)
	}

	commitMsg := "Initial commit from LokiForce Golden Path"
	_, err = w.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "LokiForce Portal",
			Email: "portal@lokiforce.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to commit files: %w", err)
	}

	headRef, err := r.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD ref: %w", err)
	}

	mainRef := plumbing.NewHashReference(plumbing.ReferenceName("refs/heads/main"), headRef.Hash())
	if err := r.Storer.SetReference(mainRef); err != nil {
		return fmt.Errorf("failed to create main branch: %w", err)
	}

	symbolicRef := plumbing.NewSymbolicReference(plumbing.HEAD, "refs/heads/main")
	if err := r.Storer.SetReference(symbolicRef); err != nil {
		return fmt.Errorf("failed to set HEAD to main: %w", err)
	}

	_ = r.Storer.RemoveReference("refs/heads/master")

	_, err = r.CreateRemote(&gitconfig.RemoteConfig{
		Name: "origin",
		URLs: []string{repoURL},
	})
	if err != nil {
		return fmt.Errorf("failed to create remote origin: %w", err)
	}

	auth := &githttp.BasicAuth{
		Username: g.owner,
		Password: g.token,
	}

	err = r.PushContext(ctx, &git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
		RefSpecs: []gitconfig.RefSpec{
			gitconfig.RefSpec("refs/heads/main:refs/heads/main"),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to push to remote repository: %w", err)
	}

	return nil
}

func copyToMemFS(src string, fs billy.Filesystem, dst string, serviceName string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := fs.MkdirAll(dstPath, 0755); err != nil {
				return err
			}
			if err := copyToMemFS(srcPath, fs, dstPath, serviceName); err != nil {
				return err
			}
		} else {
			data, err := os.ReadFile(srcPath)
			if err != nil {
				return err
			}

			content := string(data)
			content = strings.ReplaceAll(content, "lokiforce.com/scaffold/golang", serviceName)
			content = strings.ReplaceAll(content, "nodejs-scaffold", serviceName)
			content = strings.ReplaceAll(content, "# Golang Service", "# "+serviceName)
			content = strings.ReplaceAll(content, "# NodeJS Service", "# "+serviceName)

			f, err := fs.Create(dstPath)
			if err != nil {
				return err
			}
			_, err = f.Write([]byte(content))
			f.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
