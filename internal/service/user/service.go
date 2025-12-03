package user

import (
	"context"

	"github.com/anguless/mr-reviewer/internal/model"
	"github.com/anguless/mr-reviewer/internal/repository/user"
)

type userService struct {
	UserRepo user.UserRepository
}

func NewUserService(userRepo user.UserRepository) *userService {
	return &userService{
		UserRepo: userRepo,
	}
}

type UserService interface {
	UsersGetReviewGet(ctx context.Context, userID string) (*model.User, error)
	UsersSetIsActivePost(ctx context.Context, userID string, isActive bool) (*model.User, error)
	GetUser(ctx context.Context, userID string) (*model.User, error)
	GetActiveUsersByTeam(ctx context.Context, teamName string) ([]model.User, error)
	GetAssignedPRs(ctx context.Context, userID string) ([]model.PullRequestShort, error)
}
