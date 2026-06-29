package application

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
	CreateTeam(input CreateTeamInput) (CreateTeamOutput, error)
	GetTeamByID(id string) (TeamOutput, error)
	ListOrgTeams(orgID string) ([]TeamOutput, error)
}
