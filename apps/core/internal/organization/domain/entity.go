package domain

import "errors"

type Organization struct {
	ID          string
	Name        string
	Description string
	OwnerID     string
}

func NewOrganization(id, name, description, ownerID string) (*Organization, error) {
	if name == "" {
		return nil, errors.New("organization name cannot be empty")
	}
	if ownerID == "" {
		return nil, errors.New("owner ID cannot be empty")
	}
	return &Organization{
		ID:          id,
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
	}, nil
}
