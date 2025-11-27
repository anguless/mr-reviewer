package stat

import (
	"github.com/anguless/mr-reviewer/internal/db"
)

type statRepository struct {
	db db.Database
}

func NewStatRepository(db db.Database) *statRepository {
	return &statRepository{
		db: db,
	}
}
