package domain

import "context"

type ProjectRepository interface {
	Create(ctx context.Context, project *Project) error
	GetByID(ctx context.Context, id string) (*Project, error)
	ListByOrg(ctx context.Context, orgID string) ([]*Project, error)
	Delete(ctx context.Context, id string) error
}
