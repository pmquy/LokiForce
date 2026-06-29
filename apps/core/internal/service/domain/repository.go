package domain

import "context"

type ServiceRepository interface {
	Create(ctx context.Context, service *Service) error
	GetByID(ctx context.Context, id string) (*Service, error)
	ListByProject(ctx context.Context, projectID string) ([]*Service, error)
	Update(ctx context.Context, service *Service) error
	Delete(ctx context.Context, id string) error
}
