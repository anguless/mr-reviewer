package team

import (
	"context"

	"github.com/anguless/mr-reviewer/internal/db"
	"github.com/anguless/mr-reviewer/internal/model"
)

type teamRepository struct {
	db db.Database
}

func NewTeamRepository(db db.Database) *teamRepository {
	return &teamRepository{
		db: db,
	}
}

type TeamRepository interface {
	Create(ctx context.Context, team *model.Team) (*model.Team, error)
	GetByName(ctx context.Context, teamName string) (*model.Team, error)
	TeamExists(ctx context.Context, teamName string) (bool, error)
	GetTeamNameByUserID(ctx context.Context, userID string) (string, error)
}
