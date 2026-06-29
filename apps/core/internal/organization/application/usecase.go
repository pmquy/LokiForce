package application

import (
	"github.com/google/uuid"
	"lokiforce.com/apps/core/internal/organization/domain"
)

type orgUsecaseImpl struct {
	repository domain.OrganizationRepository
}

func NewOrgUsecase(repo domain.OrganizationRepository) OrgUsecase {
	return &orgUsecaseImpl{repository: repo}
}

func (u *orgUsecaseImpl) CreateOrg(input CreateOrgInput) (CreateOrgOutput, error) {
	id := uuid.NewString()
	org, err := domain.NewOrganization(id, input.Name, input.Description, input.OwnerID)
	if err != nil {
		return CreateOrgOutput{}, err
	}

	if err := u.repository.Create(org); err != nil {
		return CreateOrgOutput{}, err
	}

	return CreateOrgOutput{OrgID: org.ID}, nil
}

func (u *orgUsecaseImpl) GetOrgByID(id string) (OrgOutput, error) {
	org, err := u.repository.GetByID(id)
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

func (u *orgUsecaseImpl) ListUserOrgs(ownerID string) ([]OrgOutput, error) {
	orgs, err := u.repository.ListByOwner(ownerID)
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
