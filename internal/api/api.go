package api

import (
	"github.com/anguless/mr-reviewer/internal/api/pr"
	"github.com/anguless/mr-reviewer/internal/api/team"
	"github.com/anguless/mr-reviewer/internal/api/user"
	"github.com/anguless/mr-reviewer/internal/service"
)

type MrHandler struct {
	user.UserHandler
	team.TeamHandler
	pr.PRHandler
}

func NewMrHandler(service *service.Service) *MrHandler {
	return &MrHandler{
		UserHandler: user.NewUserHandler(service.UserService, service.PRService),
		TeamHandler: team.NewTeamHandler(service.TeamService),
		PRHandler:   pr.NewPRHandler(service.PRService),
	}
}

type Handler interface {
	pr.PRHandler
	team.TeamHandler
	user.UserHandler
}
