package user

import (
	"context"
	"errors"

	"github.com/anguless/mr-reviewer/internal/model"
)

func (s *userService) UsersGetReviewGet(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.UserRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("user is not active")
	}

	return user, nil
}

func (s *userService) UsersSetIsActivePost(ctx context.Context, userID string, isActive bool) (*model.User, error) {
	user, err := s.UserRepo.UpdateIsActive(ctx, userID, isActive)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUser(ctx context.Context, userID string) (*model.User, error) {
	return s.UserRepo.GetByID(ctx, userID)
}

func (s *userService) GetActiveUsersByTeam(ctx context.Context, teamName string) ([]model.User, error) {
	return s.UserRepo.GetActiveByTeam(ctx, teamName)
}

func (s *userService) GetAssignedPRs(ctx context.Context, userID string) ([]model.PullRequestShort, error) {
	return s.UserRepo.GetAssignedPRs(ctx, userID)
}
