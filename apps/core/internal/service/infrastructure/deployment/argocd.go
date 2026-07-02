package deployment

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"text/template"

	"lokiforce.com/apps/core/internal/config"
	"lokiforce.com/apps/core/internal/service/application"
	gitclient "lokiforce.com/apps/core/pkg/git"
)

type ArgoCDDeployment struct {
	client     gitclient.GitClient
	gitopsRepo string
	owner      string
}

func NewArgoCDDeployment(cfg *config.Config, client gitclient.GitClient) application.DeploymentControl {
	return &ArgoCDDeployment{
		client:     client,
		gitopsRepo: cfg.GitHub.GitOpsRepo,
		owner:      cfg.GitHub.Owner,
	}
}

func (a *ArgoCDDeployment) RegisterDeployment(ctx context.Context, config application.DeploymentConfig) (string, error) {

	data := map[string]any{
		"ServiceName": config.ServiceName,
		"Owner":       a.owner,
		"GitOpsRepo":  a.gitopsRepo,
		"Namespace":   config.Namespace,
	}

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

		slog.Info("Committing manifest via client", "repo", a.gitopsRepo, "path", destPath)
		message := fmt.Sprintf("Add GitOps manifest: %s", destPath)
		if err := a.client.WriteFile(ctx, a.gitopsRepo, destPath, buf.String(), message); err != nil {
			return "", fmt.Errorf("failed to commit %s: %w", destPath, err)
		}
	}

	remoteURL := fmt.Sprintf("https://github.com/%s/%s/blob/main/apps/%s-argocd.yaml", a.owner, a.gitopsRepo, config.ServiceName)
	slog.Info("Argo CD and K8s manifests committed to remote GitOps repo successfully", "url", remoteURL)
	return remoteURL, nil
}

func (a *ArgoCDDeployment) DeleteDeployment(ctx context.Context, serviceName string) error {

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

	for _, relPath := range filesToDelete {
		slog.Info("Deleting remote manifest via client", "repo", a.gitopsRepo, "path", relPath)
		message := fmt.Sprintf("Delete GitOps manifest: %s", relPath)
		if err := a.client.DeleteFile(ctx, a.gitopsRepo, relPath, message); err != nil {
			return fmt.Errorf("failed to delete remote file %s: %w", relPath, err)
		}
	}

	slog.Info("Argo CD and K8s manifests deleted from remote GitOps repo successfully", "serviceName", serviceName)
	return nil
}

func (a *ArgoCDDeployment) CreateAccessPolicy(ctx context.Context, policyID, clientID, targetID, targetPort, namespace string) (string, error) {

	data := map[string]any{
		"PolicyID":   policyID,
		"ClientID":   clientID,
		"TargetID":   targetID,
		"TargetPort": targetPort,
		"Namespace":  namespace,
	}

	tmpl, err := template.New("network-policy").Parse(networkPolicyTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse network policy template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute network policy template: %w", err)
	}

	relPath := fmt.Sprintf("policies/policy-%s.yaml", policyID)
	message := fmt.Sprintf("Add NetworkPolicy: %s", relPath)

	slog.Info("Committing NetworkPolicy via client", "repo", a.gitopsRepo, "path", relPath)
	if err := a.client.WriteFile(ctx, a.gitopsRepo, relPath, buf.String(), message); err != nil {
		return "", fmt.Errorf("failed to write policy file: %w", err)
	}

	remoteURL := fmt.Sprintf("https://github.com/%s/%s/blob/main/%s", a.owner, a.gitopsRepo, relPath)
	return remoteURL, nil
}

func (a *ArgoCDDeployment) DeleteAccessPolicy(ctx context.Context, policyID string) error {
	relPath := fmt.Sprintf("policies/policy-%s.yaml", policyID)
	message := fmt.Sprintf("Delete NetworkPolicy: %s", relPath)

	slog.Info("Deleting NetworkPolicy via client", "repo", a.gitopsRepo, "path", relPath)
	if err := a.client.DeleteFile(ctx, a.gitopsRepo, relPath, message); err != nil {
		return fmt.Errorf("failed to delete policy file: %w", err)
	}

	return nil
}
