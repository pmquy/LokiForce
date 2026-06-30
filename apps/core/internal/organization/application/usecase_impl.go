package application

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"lokiforce.com/apps/core/internal/organization/domain"
	"lokiforce.com/apps/core/pkg/mq"
)

type orgUsecaseImpl struct {
	repository domain.OrganizationRepository
	mq         mq.MessageQueue
}

func NewOrgUsecase(repo domain.OrganizationRepository, msgQueue mq.MessageQueue) OrgUsecase {
	u := &orgUsecaseImpl{
		repository: repo,
		mq:         msgQueue,
	}

	msgQueue.Subscribe("user.registered", func(payload any) {
		event, ok := payload.(mq.UserRegisteredEvent)
		if !ok {
			slog.Error("Invalid payload type for user.registered event")
			return
		}

		slog.Info("Event user.registered received asynchronously. Creating default organization...", "userID", event.UserID, "username", event.Username)

		ctx := context.Background()
		_, err := u.CreateOrg(ctx, CreateOrgInput{
			Name:        event.Username + "'s Org",
			Description: "Default organization created automatically upon registration",
			OwnerID:     event.UserID,
		})
		if err != nil {
			slog.Error("Failed to automatically create default organization", "error", err, "userID", event.UserID)
		}
	})

	return u
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
