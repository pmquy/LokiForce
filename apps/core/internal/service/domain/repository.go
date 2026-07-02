package domain

import "context"

type ServiceRepository interface {
	Create(ctx context.Context, service *Service) error
	GetByID(ctx context.Context, id string) (*Service, error)
	ListByProject(ctx context.Context, projectID string) ([]*Service, error)
	Update(ctx context.Context, service *Service) error
	Delete(ctx context.Context, id string) error
	CreateAccessPolicy(ctx context.Context, policy *AccessPolicy) error
	DeleteAccessPolicy(ctx context.Context, policyID string) error
	ListAccessPoliciesByTarget(ctx context.Context, targetID string) ([]*AccessPolicy, error)
	GetAccessPolicyByID(ctx context.Context, id string) (*AccessPolicy, error)
}
