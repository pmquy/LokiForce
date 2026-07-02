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

type AccessPolicy struct {
	ID         string
	ClientID   string
	TargetID   string
	TargetPort string
	ProjectID  string
}

func NewAccessPolicy(id, clientID, targetID, targetPort, projectID string) (*AccessPolicy, error) {
	if clientID == "" {
		return nil, errors.New("client ID cannot be empty")
	}
	if targetID == "" {
		return nil, errors.New("target ID cannot be empty")
	}
	if targetPort == "" {
		return nil, errors.New("target port cannot be empty")
	}
	if projectID == "" {
		return nil, errors.New("project ID cannot be empty")
	}
	return &AccessPolicy{
		ID:         id,
		ClientID:   clientID,
		TargetID:   targetID,
		TargetPort: targetPort,
		ProjectID:  projectID,
	}, nil
}
