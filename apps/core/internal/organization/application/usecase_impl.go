package application

import (
	"context"

	"github.com/google/uuid"
	"lokiforce.com/apps/core/internal/organization/domain"
)

type orgUsecaseImpl struct {
	repository domain.OrganizationRepository
}

func NewOrgUsecase(repo domain.OrganizationRepository) OrgUsecase {
	return &orgUsecaseImpl{repository: repo}
}

func (u *orgUsecaseImpl) CreateOrg(ctx context.Context, input CreateOrgInput) (CreateOrgOutput, error) {
	id := uuid.NewString()
	org, err := domain.NewOrganization(id, input.Name, input.Description, input.OwnerID)
	if err != nil {
		return CreateOrgOutput{}, err
	}

	if err := u.repository.Create(ctx, org); err != nil {
		return CreateOrgOutput{}, err
	}

	return CreateOrgOutput{OrgID: org.ID}, nil
}

func (u *orgUsecaseImpl) GetOrgByID(ctx context.Context, id string) (OrgOutput, error) {
	org, err := u.repository.GetByID(ctx, id)
	if err != nil {
		return OrgOutput{}, err
	}

	return OrgOutput{
		ID:          org.ID,
		Name:        org.Name,
		Description: org.Description,
		OwnerID:     org.OwnerID,
	}, nil
}

func (u *orgUsecaseImpl) ListUserOrgs(ctx context.Context, ownerID string) ([]OrgOutput, error) {
	orgs, err := u.repository.ListByOwner(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	outputs := make([]OrgOutput, len(orgs))
	for i, org := range orgs {
		outputs[i] = OrgOutput{
			ID:          org.ID,
			Name:        org.Name,
			Description: org.Description,
			OwnerID:     org.OwnerID,
		}
	}
	return outputs, nil
}
