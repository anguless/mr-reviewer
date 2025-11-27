package user

import (
	"github.com/anguless/mr-reviewer/internal/db"
)

type userRepository struct {
	db db.Database
}

func NewUserRepository(db db.Database) *userRepository {
	return &userRepository{
		db: db,
	}
}
