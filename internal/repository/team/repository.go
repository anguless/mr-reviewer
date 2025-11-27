package team

import (
	"github.com/anguless/mr-reviewer/internal/db"
)

type teamRepository struct {
	db db.Database
}

func NewTeamRepository(db db.Database) *teamRepository {
	return &teamRepository{
		db: db,
	}
}
