package user

import (
	"context"
	"errors"

	"github.com/anguless/reviewer/internal/model"
	"github.com/google/uuid"
)

func (s *UserService) CreateUser(ctx context.Context, userC *model.User) (*model.User, error) {
	user := &model.User{
		ID:       uuid.New(),
		Username: userC.Username,
		TeamID:   userC.TeamID,
		IsActive: userC.IsActive,
	}
	err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	return s.userRepo.UpdateUser(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.DeleteUser(ctx, id)
}

func (s *UserService) GetAssignedPRs(ctx context.Context, userID uuid.UUID) ([]model.PullRequest, error) {
	return s.userRepo.GetPRsByReviewer(ctx, userID)
}
