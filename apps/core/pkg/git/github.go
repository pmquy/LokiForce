package git

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
)

type GitHubClient struct {
	token string
	owner string
}

func NewGitHubClient(token, owner string) GitClient {
	return &GitHubClient{
		token: token,
		owner: owner,
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

func (c *GitHubClient) CreateRepository(ctx context.Context, name, description string, private bool) (string, error) {
	reqBody := createRepoRequest{
		Name:        name,
		Description: description,
		Private:     private,
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

	c.setHeaders(req)
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

func (c *GitHubClient) DeleteRepository(ctx context.Context, name string) error {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s", c.owner, name)
	req, err := http.NewRequestWithContext(ctx, "DELETE", apiURL, nil)
	if err != nil {
		return err
	}

	c.setHeaders(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("github delete api failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (c *GitHubClient) GetFileSHA(ctx context.Context, repo, path string) (string, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", c.owner, repo, path)
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return "", err
	}

	c.setHeaders(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("github get content api failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	type fileInfoResponse struct {
		SHA string `json:"sha"`
	}
	var fileInfo fileInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&fileInfo); err != nil {
		return "", err
	}
	return fileInfo.SHA, nil
}

func (c *GitHubClient) WriteFile(ctx context.Context, repo, path, content, message string) error {
	base64Content := base64.StdEncoding.EncodeToString([]byte(content))

	type gitHubContentWithSHARequest struct {
		Message string `json:"message"`
		Content string `json:"content"`
		SHA     string `json:"sha,omitempty"`
	}

	reqBody := gitHubContentWithSHARequest{
		Message: message,
		Content: base64Content,
	}

	sha, _ := c.GetFileSHA(ctx, repo, path)
	if sha != "" {
		reqBody.SHA = sha
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", c.owner, repo, path)
	req, err := http.NewRequestWithContext(ctx, "PUT", apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to write %s to github: status %d, response %s", path, resp.StatusCode, string(respBody))
	}

	return nil
}

func (c *GitHubClient) DeleteFile(ctx context.Context, repo, path, message string) error {
	sha, err := c.GetFileSHA(ctx, repo, path)
	if err != nil {
		return err
	}
	if sha == "" {
		return nil
	}

	type deleteContentRequest struct {
		Message string `json:"message"`
		SHA     string `json:"sha"`
	}
	reqBody := deleteContentRequest{
		Message: message,
		SHA:     sha,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", c.owner, repo, path)
	req, err := http.NewRequestWithContext(ctx, "DELETE", apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete file %s: status %d, response %s", path, resp.StatusCode, string(respBody))
	}

	return nil
}

func (c *GitHubClient) PushFiles(ctx context.Context, repoURL string, files map[string][]byte, commitMessage string) error {
	fs := memfs.New()
	store := memory.NewStorage()

	slog.Info("Initializing local in-memory repository", "url", repoURL)
	repo, err := git.Init(store, fs)
	if err != nil {
		return fmt.Errorf("failed to init repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	for path, content := range files {
		dir := filepath.Dir(path)
		if dir != "." && dir != "/" {
			if err := fs.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}
		}

		f, err := fs.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", path, err)
		}
		if _, err := f.Write(content); err != nil {
			f.Close()
			return fmt.Errorf("failed to write file %s: %w", path, err)
		}
		f.Close()
	}

	_, err = worktree.Add(".")
	if err != nil {
		return fmt.Errorf("failed to add files to git worktree: %w", err)
	}

	_, err = worktree.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Git Client Bot",
			Email: "bot@gitclient.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to commit files: %w", err)
	}

	headRef, err := repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD ref: %w", err)
	}

	mainRef := plumbing.NewHashReference(plumbing.ReferenceName("refs/heads/main"), headRef.Hash())
	if err := repo.Storer.SetReference(mainRef); err != nil {
		return fmt.Errorf("failed to create main branch: %w", err)
	}

	symbolicRef := plumbing.NewSymbolicReference(plumbing.HEAD, "refs/heads/main")
	if err := repo.Storer.SetReference(symbolicRef); err != nil {
		return fmt.Errorf("failed to set HEAD to main: %w", err)
	}

	_ = repo.Storer.RemoveReference("refs/heads/master")

	_, err = repo.CreateRemote(&gitconfig.RemoteConfig{
		Name: "origin",
		URLs: []string{repoURL},
	})
	if err != nil {
		return fmt.Errorf("failed to create remote origin: %w", err)
	}

	err = repo.PushContext(ctx, &git.PushOptions{
		RemoteName: "origin",
		Auth:       &githttp.BasicAuth{Username: "git", Password: c.token},
		RefSpecs:   []gitconfig.RefSpec{"refs/heads/main:refs/heads/main"},
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("failed to push to remote repository: %w", err)
	}

	slog.Info("Files pushed to remote repository successfully", "url", repoURL)
	return nil
}

func (c *GitHubClient) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
}
