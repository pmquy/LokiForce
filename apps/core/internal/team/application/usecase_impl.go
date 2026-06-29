package application

import (
	"context"

	"github.com/google/uuid"
	"lokiforce.com/apps/core/internal/team/domain"
)

type teamUsecaseImpl struct {
	repository domain.TeamRepository
}

func NewTeamUsecase(repo domain.TeamRepository) TeamUsecase {
	return &teamUsecaseImpl{repository: repo}
}

func (u *teamUsecaseImpl) CreateTeam(ctx context.Context, input CreateTeamInput) (CreateTeamOutput, error) {
	id := uuid.NewString()
	team, err := domain.NewTeam(id, input.Name, input.Description, input.OrgID)
	if err != nil {
		return CreateTeamOutput{}, err
	}

	if err := u.repository.Create(ctx, team); err != nil {
		return CreateTeamOutput{}, err
	}

	return CreateTeamOutput{TeamID: team.ID}, nil
}

func (u *teamUsecaseImpl) GetTeamByID(ctx context.Context, id string) (TeamOutput, error) {
	team, err := u.repository.GetByID(ctx, id)
	if err != nil {
		return TeamOutput{}, err
	}

	return TeamOutput{
		ID:          team.ID,
		Name:        team.Name,
		Description: team.Description,
		OrgID:       team.OrgID,
	}, nil
}

func (u *teamUsecaseImpl) ListOrgTeams(ctx context.Context, orgID string) ([]TeamOutput, error) {
	teams, err := u.repository.ListByOrg(ctx, orgID)
	if err != nil {
		return nil, err
	}

	outputs := make([]TeamOutput, len(teams))
	for i, team := range teams {
		outputs[i] = TeamOutput{
			ID:          team.ID,
			Name:        team.Name,
			Description: team.Description,
			OrgID:       team.OrgID,
		}
	}
	return outputs, nil
}
