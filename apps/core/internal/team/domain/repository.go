package domain

type TeamRepository interface {
	Create(team *Team) error
	GetByID(id string) (*Team, error)
	ListByOrg(orgID string) ([]*Team, error)
	Delete(id string) error
}
