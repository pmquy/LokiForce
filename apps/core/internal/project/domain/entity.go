package domain

import "errors"

type Project struct {
	ID          string
	Name        string
	Description string
	OrgID       string
}

func NewProject(id, name, description, orgID string) (*Project, error) {
	if name == "" {
		return nil, errors.New("project name cannot be empty")
	}
	if orgID == "" {
		return nil, errors.New("organization ID cannot be empty")
	}
	return &Project{
		ID:          id,
		Name:        name,
		Description: description,
		OrgID:       orgID,
	}, nil
}
