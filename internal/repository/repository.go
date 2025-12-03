package repository

import (
	"github.com/anguless/mr-reviewer/internal/db"
	"github.com/anguless/mr-reviewer/internal/repository/pr"
	"github.com/anguless/mr-reviewer/internal/repository/team"
	"github.com/anguless/mr-reviewer/internal/repository/user"
)

type Repo struct {
	UserRepo user.UserRepository
	TeamRepo team.TeamRepository
	PRRepo   pr.PRRepository
}

func NewRepository(db db.Database) *Repo {
	return &Repo{
		UserRepo: user.NewUserRepository(db),
		TeamRepo: team.NewTeamRepository(db),
		PRRepo:   pr.NewPRRepository(db),
	}
}
