package team

import (
	"context"

	"github.com/anguless/mr-reviewer/internal/model"
	"github.com/anguless/mr-reviewer/internal/repository/team"
)

type teamService struct {
	TeamRepository team.TeamRepository
}

func NewTeamService(teamRepo team.TeamRepository) *teamService {
	return &teamService{
		TeamRepository: teamRepo,
	}
}

type TeamService interface {
	TeamAddPost(ctx context.Context, team *model.Team) (*model.Team, error)
	TeamGetGet(ctx context.Context, teamName string) (*model.Team, error)
	GetTeamNameByUser(ctx context.Context, userID string) (string, error)
}
