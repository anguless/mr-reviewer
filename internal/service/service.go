package service

import (
	"github.com/anguless/mr-reviewer/internal/repository"
	"github.com/anguless/mr-reviewer/internal/service/pr"
	"github.com/anguless/mr-reviewer/internal/service/team"
	"github.com/anguless/mr-reviewer/internal/service/user"
	service "github.com/anguless/mr-reviewer/internal/service/util"
)

type Service struct {
	UserService user.UserService
	TeamService team.TeamService
	PRService   pr.PRService
}

func NewService(repo *repository.Repo, rnd service.RandomGenerator) *Service {
	userService := user.NewUserService(repo.UserRepo)
	teamService := team.NewTeamService(repo.TeamRepo)
	prService := pr.NewPRService(repo.PRRepo, repo.UserRepo, repo.TeamRepo, rnd)

	return &Service{
		UserService: userService,
		TeamService: teamService,
		PRService:   prService,
	}
}
