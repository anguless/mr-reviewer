package pr

import "github.com/anguless/reviewer/internal/db"

type prRepository struct {
	db db.Database
}

func NewPRRepository(db db.Database) *prRepository {
	return &prRepository{
		db: db,
	}
}
