package domain

type ProjectRepository interface {
	Create(project *Project) error
	GetByID(id string) (*Project, error)
	ListByOrg(orgID string) ([]*Project, error)
	Delete(id string) error
}
