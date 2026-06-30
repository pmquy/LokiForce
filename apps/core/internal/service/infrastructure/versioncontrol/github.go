package versioncontrol

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

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

func (g *GitHubVersionControl) PushFiles(ctx context.Context, repoURL string, localDir string) error {

	if g.token == "" || g.token == "mock_token" {
		return nil
	}

	authURL := repoURL
	if strings.HasPrefix(repoURL, "https://") {
		authURL = "https://" + g.token + "@" + strings.TrimPrefix(repoURL, "https://")
	}

	runCmd := func(name string, args ...string) error {
		cmd := exec.CommandContext(ctx, name, args...)
		cmd.Dir = localDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	if err := runCmd("git", "init"); err != nil {
		return fmt.Errorf("git init failed: %w", err)
	}

	_ = runCmd("git", "config", "user.name", "LokiForce Portal")
	_ = runCmd("git", "config", "user.email", "portal@lokiforce.com")

	if err := runCmd("git", "add", "."); err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}
	if err := runCmd("git", "commit", "-m", "Initial commit from LokiForce Golden Path"); err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}
	if err := runCmd("git", "branch", "-M", "main"); err != nil {
		return fmt.Errorf("git branch rename failed: %w", err)
	}

	_ = runCmd("git", "remote", "remove", "origin")
	if err := runCmd("git", "remote", "add", "origin", authURL); err != nil {
		return fmt.Errorf("git remote add failed: %w", err)
	}

	if err := runCmd("git", "push", "-u", "origin", "main"); err != nil {
		return fmt.Errorf("git push failed: %w", err)
	}

	return nil
}
