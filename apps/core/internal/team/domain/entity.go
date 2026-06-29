package domain

import "errors"

type Team struct {
	ID          string
	Name        string
	Description string
	OrgID       string
}

func NewTeam(id, name, description, orgID string) (*Team, error) {
	if name == "" {
		return nil, errors.New("team name cannot be empty")
	}
	if orgID == "" {
		return nil, errors.New("organization ID cannot be empty")
	}
	return &Team{
		ID:          id,
		Name:        name,
		Description: description,
		OrgID:       orgID,
	}, nil
}
