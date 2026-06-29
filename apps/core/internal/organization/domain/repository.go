package domain

import "context"

type OrganizationRepository interface {
	Create(ctx context.Context, org *Organization) error
	GetByID(ctx context.Context, id string) (*Organization, error)
	ListByOwner(ctx context.Context, ownerID string) ([]*Organization, error)
	Delete(ctx context.Context, id string) error
}
