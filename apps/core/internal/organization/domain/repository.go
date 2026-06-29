package domain

type OrganizationRepository interface {
	Create(org *Organization) error
	GetByID(id string) (*Organization, error)
	ListByOwner(ownerID string) ([]*Organization, error)
	Delete(id string) error
}
