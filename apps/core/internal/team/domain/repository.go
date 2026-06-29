package domain

import "context"

type TeamRepository interface {
	Create(ctx context.Context, team *Team) error
	GetByID(ctx context.Context, id string) (*Team, error)
	ListByOrg(ctx context.Context, orgID string) ([]*Team, error)
	Delete(ctx context.Context, id string) error
}
