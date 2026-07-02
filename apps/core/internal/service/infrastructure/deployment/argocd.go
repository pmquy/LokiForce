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

const deploymentTemplate = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.ServiceName}}
  namespace: production
  labels:
    app: {{.ServiceName}}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{.ServiceName}}
  template:
    metadata:
      labels:
        app: {{.ServiceName}}
    spec:
      containers:
        - name: app
          image: ghcr.io/{{.Owner}}/{{.ServiceName}}:latest
          ports:
            - containerPort: 8080
          resources:
            limits:
              cpu: "500m"
              memory: "512Mi"
            requests:
              cpu: "100m"
              memory: "128Mi"
`

const serviceTemplate = `apiVersion: v1
kind: Service
metadata:
  name: {{.ServiceName}}
  namespace: production
  labels:
    app: {{.ServiceName}}
spec:
  ports:
    - port: 80
      targetPort: 8080
      protocol: TCP
  selector:
    app: {{.ServiceName}}
  type: ClusterIP
`

const kustomizationTemplate = `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - deployment.yaml
  - service.yaml
images:
  - name: ghcr.io/{{.Owner}}/{{.ServiceName}}
    newTag: latest
`

const argocdTemplate = `apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: {{.ServiceName}}
  namespace: argocd
  annotations:
    argocd-image-updater.argoproj.io/image-list: app=ghcr.io/{{.Owner}}/{{.ServiceName}}
    argocd-image-updater.argoproj.io/app.update-strategy: latest
    argocd-image-updater.argoproj.io/write-back-method: git:kustomize
    argocd-image-updater.argoproj.io/write-back-path: manifests/{{.ServiceName}}
  finalizers:
    - resources-finalizer.argocd.argoproj.io
spec:
  project: default
  source:
    repoURL: 'https://github.com/{{.Owner}}/{{.GitOpsRepo}}.git'
    targetRevision: HEAD
    path: manifests/{{.ServiceName}}
  destination:
    server: 'https://kubernetes.default.svc'
    namespace: {{.Namespace}}
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
`

func (a *ArgoCDDeployment) RegisterDeployment(ctx context.Context, config application.DeploymentConfig) (string, error) {
	if a.token == "" {
		return "", fmt.Errorf("github token is empty: authorization failed")
	}

	renderStrTemplate := func(name string, tmplStr string) (string, error) {
		tmpl, err := template.New(name).Parse(tmplStr)
		if err != nil {
			return "", err
		}
		data := map[string]any{
			"ServiceName": config.ServiceName,
			"Owner":       a.owner,
			"GitOpsRepo":  a.gitopsRepo,
			"Namespace":   config.Namespace,
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return "", err
		}
		return buf.String(), nil
	}

	argocdYaml, err := renderStrTemplate("argocd", argocdTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to render argocd template: %w", err)
	}

	deploymentYaml, err := renderStrTemplate("deployment", deploymentTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to render deployment template: %w", err)
	}

	serviceYaml, err := renderStrTemplate("service", serviceTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to render service template: %w", err)
	}

	kustomizationYaml, err := renderStrTemplate("kustomization", kustomizationTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to render kustomization template: %w", err)
	}

	filesToWrite := map[string]string{
		fmt.Sprintf("apps/%s-argocd.yaml", config.ServiceName):                 argocdYaml,
		fmt.Sprintf("manifests/%s/deployment.yaml", config.ServiceName):         deploymentYaml,
		fmt.Sprintf("manifests/%s/service.yaml", config.ServiceName):            serviceYaml,
		fmt.Sprintf("manifests/%s/kustomization.yaml", config.ServiceName):      kustomizationYaml,
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

	filesToDelete := []string{
		fmt.Sprintf("apps/%s-argocd.yaml", serviceName),
		fmt.Sprintf("manifests/%s/deployment.yaml", serviceName),
		fmt.Sprintf("manifests/%s/service.yaml", serviceName),
		fmt.Sprintf("manifests/%s/kustomization.yaml", serviceName),
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
