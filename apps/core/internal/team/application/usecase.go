package application

import "context"

type CreateTeamInput struct {
	Name        string
	Description string
	OrgID       string
}

type CreateTeamOutput struct {
	TeamID string
}

type TeamOutput struct {
	ID          string
	Name        string
	Description string
	OrgID       string
}

type TeamUsecase interface {
	CreateTeam(ctx context.Context, input CreateTeamInput) (CreateTeamOutput, error)
	GetTeamByID(ctx context.Context, id string) (TeamOutput, error)
	ListOrgTeams(ctx context.Context, orgID string) ([]TeamOutput, error)
}
