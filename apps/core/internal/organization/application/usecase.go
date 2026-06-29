package application

import "context"

type CreateOrgInput struct {
	Name        string
	Description string
	OwnerID     string
}

type CreateOrgOutput struct {
	OrgID string
}

type OrgOutput struct {
	ID          string
	Name        string
	Description string
	OwnerID     string
}

type OrgUsecase interface {
	CreateOrg(ctx context.Context, input CreateOrgInput) (CreateOrgOutput, error)
	GetOrgByID(ctx context.Context, id string) (OrgOutput, error)
	ListUserOrgs(ctx context.Context, ownerID string) ([]OrgOutput, error)
}
