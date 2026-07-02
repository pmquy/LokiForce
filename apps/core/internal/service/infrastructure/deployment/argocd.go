package deployment

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"text/template"

	"lokiforce.com/apps/core/internal/config"
	"lokiforce.com/apps/core/internal/service/application"
)

type ArgoCDDeployment struct {
	token      string
	owner      string
	gitopsRepo string
}

func NewArgoCDDeployment(cfg *config.Config) application.DeploymentControl {
	return &ArgoCDDeployment{
		token:      cfg.GitHub.Token,
		owner:      cfg.GitHub.Owner,
		gitopsRepo: cfg.GitHub.GitOpsRepo,
	}
}

type gitHubContentRequest struct {
	Message string `json:"message"`
	Content string `json:"content"`
}

func (a *ArgoCDDeployment) RegisterDeployment(ctx context.Context, config application.DeploymentConfig) (string, error) {
	if a.token == "" {
		return "", fmt.Errorf("github token is empty: authorization failed")
	}

	data := map[string]interface{}{
		"ServiceName": config.ServiceName,
		"Owner":       a.owner,
		"GitOpsRepo":  a.gitopsRepo,
		"Namespace":   config.Namespace,
	}

	filesToWrite := make(map[string]string)
	for fileName, rawTemplate := range TemplatesMap {
		tmpl, err := template.New(fileName).Parse(rawTemplate)
		if err != nil {
			return "", fmt.Errorf("failed to parse template %s: %w", fileName, err)
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return "", fmt.Errorf("failed to execute template %s: %w", fileName, err)
		}

		var destPath string
		if fileName == "argocd.yaml" {
			destPath = fmt.Sprintf("apps/%s-argocd.yaml", config.ServiceName)
		} else {
			destPath = fmt.Sprintf("manifests/%s/%s", config.ServiceName, fileName)
		}

		filesToWrite[destPath] = buf.String()
	}

	writeRemoteFile := func(path string, content string) error {
		base64Content := base64.StdEncoding.EncodeToString([]byte(content))
		reqBody := gitHubContentRequest{
			Message: fmt.Sprintf("Add GitOps manifest: %s", path),
			Content: base64Content,
		}
		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			return err
		}

		apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", a.owner, a.gitopsRepo, path)
		req, err := http.NewRequestWithContext(ctx, "PUT", apiURL, bytes.NewBuffer(bodyBytes))
		if err != nil {
			return err
		}

		req.Header.Set("Authorization", "Bearer "+a.token)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
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

	for relPath, content := range filesToWrite {
		slog.Info("Committing manifest to remote GitOps repo", "repo", a.gitopsRepo, "path", relPath)
		if err := writeRemoteFile(relPath, content); err != nil {
			return "", fmt.Errorf("failed to commit %s: %w", relPath, err)
		}
	}

	remoteURL := fmt.Sprintf("https://github.com/%s/%s/blob/main/apps/%s-argocd.yaml", a.owner, a.gitopsRepo, config.ServiceName)
	slog.Info("Argo CD and K8s manifests committed to remote GitOps repo successfully", "url", remoteURL)
	return remoteURL, nil
}

func (a *ArgoCDDeployment) DeleteDeployment(ctx context.Context, serviceName string) error {
	if a.token == "" {
		return fmt.Errorf("github token is empty: authorization failed")
	}

	var filesToDelete []string
	for fileName := range TemplatesMap {
		var relPath string
		if fileName == "argocd.yaml" {
			relPath = fmt.Sprintf("apps/%s-argocd.yaml", serviceName)
		} else {
			relPath = fmt.Sprintf("manifests/%s/%s", serviceName, fileName)
		}
		filesToDelete = append(filesToDelete, relPath)
	}

	deleteRemoteFile := func(path string) error {

		apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", a.owner, a.gitopsRepo, path)
		req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+a.token)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			return nil
		}
		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to get file %s info: status %d, response %s", path, resp.StatusCode, string(respBody))
		}

		type fileInfoResponse struct {
			SHA string `json:"sha"`
		}
		var fileInfo fileInfoResponse
		if err := json.NewDecoder(resp.Body).Decode(&fileInfo); err != nil {
			return err
		}

		type deleteContentRequest struct {
			Message string `json:"message"`
			SHA     string `json:"sha"`
		}
		reqBody := deleteContentRequest{
			Message: fmt.Sprintf("Delete GitOps manifest: %s", path),
			SHA:     fileInfo.SHA,
		}
		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			return err
		}

		reqDel, err := http.NewRequestWithContext(ctx, "DELETE", apiURL, bytes.NewBuffer(bodyBytes))
		if err != nil {
			return err
		}
		reqDel.Header.Set("Authorization", "Bearer "+a.token)
		reqDel.Header.Set("Accept", "application/vnd.github+json")
		reqDel.Header.Set("X-GitHub-Api-Version", "2022-11-28")
		reqDel.Header.Set("Content-Type", "application/json")

		respDel, err := client.Do(reqDel)
		if err != nil {
			return err
		}
		defer respDel.Body.Close()

		if respDel.StatusCode != http.StatusOK && respDel.StatusCode != http.StatusNotFound {
			respBody, _ := io.ReadAll(respDel.Body)
			return fmt.Errorf("failed to delete file %s from GitOps repo: status %d, response %s", path, respDel.StatusCode, string(respBody))
		}
		return nil
	}

	for _, relPath := range filesToDelete {
		slog.Info("Deleting remote manifest", "repo", a.gitopsRepo, "path", relPath)
		if err := deleteRemoteFile(relPath); err != nil {
			return fmt.Errorf("failed to delete remote file %s: %w", relPath, err)
		}
	}

	slog.Info("Argo CD and K8s manifests deleted from remote GitOps repo successfully", "serviceName", serviceName)
	return nil
}

func (a *ArgoCDDeployment) CreateAccessPolicy(ctx context.Context, clientID, targetID, targetPort, projectID string) (string, error) {
	return "", nil
}

func (a *ArgoCDDeployment) DeleteAccessPolicy(ctx context.Context, policyID string) error {
	return nil
}
