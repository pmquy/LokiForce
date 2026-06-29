package domain

import "errors"

type Service struct {
	ID          string
	Name        string
	Description string
	ProjectID   string
	TemplateID  string
	Repository  string
}

func NewService(id, name, description, projectID, templateID string) (*Service, error) {
	if name == "" {
		return nil, errors.New("service name cannot be empty")
	}
	if projectID == "" {
		return nil, errors.New("project ID cannot be empty")
	}
	if templateID == "" {
		return nil, errors.New("template ID cannot be empty")
	}
	return &Service{
		ID:          id,
		Name:        name,
		Description: description,
		ProjectID:   projectID,
		TemplateID:  templateID,
	}, nil
}
