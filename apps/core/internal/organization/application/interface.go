package application

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
	CreateOrg(input CreateOrgInput) (CreateOrgOutput, error)
	GetOrgByID(id string) (OrgOutput, error)
	ListUserOrgs(ownerID string) ([]OrgOutput, error)
}
