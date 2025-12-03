package team

import (
	"context"

	"github.com/anguless/mr-reviewer/internal/service/team"
	mr_v1 "github.com/anguless/mr-reviewer/pkg/openapi/mr/v1"
)

type teamHandler struct {
	teamService team.TeamService
}

func NewTeamHandler(teamService team.TeamService) *teamHandler {
	return &teamHandler{
		teamService: teamService,
	}
}

type TeamHandler interface {
	TeamAddPost(ctx context.Context, req *mr_v1.Team) (mr_v1.TeamAddPostRes, error)
	TeamGetGet(ctx context.Context, params mr_v1.TeamGetGetParams) (mr_v1.TeamGetGetRes, error)
}
