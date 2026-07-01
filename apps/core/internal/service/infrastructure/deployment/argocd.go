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

	yamlContent := fmt.Sprintf(`apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: %s
  namespace: argocd
  finalizers:
    - resources-finalizer.argocd.argoproj.io
spec:
  project: default
  source:
    repoURL: '%s'
    targetRevision: HEAD
    path: deploy/k8s
  destination:
    server: 'https://kubernetes.default.svc'
    namespace: %s
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
`, config.ServiceName, config.RepositoryURL, config.Namespace)

	slog.Info("Committing Argo CD manifest to remote GitOps repo", "repo", a.gitopsRepo, "file", config.ServiceName)
	base64Content := base64.StdEncoding.EncodeToString([]byte(yamlContent))
	reqBody := gitHubContentRequest{
		Message: fmt.Sprintf("Add ArgoCD Application for service %s", config.ServiceName),
		Content: base64Content,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/apps/%s-argocd.yaml", a.owner, a.gitopsRepo, config.ServiceName)
	req, err := http.NewRequestWithContext(ctx, "PUT", apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+a.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to commit manifest to gitops repo: status %d, response %s", resp.StatusCode, string(respBody))
	}

	remoteURL := fmt.Sprintf("https://github.com/%s/%s/blob/main/apps/%s-argocd.yaml", a.owner, a.gitopsRepo, config.ServiceName)
	slog.Info("Argo CD Application manifest committed to remote GitOps repo successfully", "url", remoteURL)
	return remoteURL, nil
}
