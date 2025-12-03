package user

import (
	"context"

	"github.com/anguless/mr-reviewer/internal/db"
	"github.com/anguless/mr-reviewer/internal/model"
)

type userRepository struct {
	db db.Database
}

func NewUserRepository(db db.Database) *userRepository {
	return &userRepository{
		db: db,
	}
}

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	GetByID(ctx context.Context, userID string) (*model.User, error)
	GetByTeam(ctx context.Context, teamName string) ([]model.User, error)
	UpdateIsActive(ctx context.Context, userID string, isActive bool) (*model.User, error)
	GetActiveByTeam(ctx context.Context, teamName string) ([]model.User, error)
	GetAssignedPRs(ctx context.Context, userID string) ([]model.PullRequestShort, error)
}
