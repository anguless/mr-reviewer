package team

import (
	"context"
	"errors"
	"fmt"

	"github.com/anguless/mr-reviewer/internal/model"
	"github.com/google/uuid"
)

func (s *TeamService) CreateTeam(ctx context.Context, team *model.Team) (*model.Team, error) {
	existing, err := s.teamRepo.GetByName(ctx, team.Name)

	if err != nil {
		return nil, fmt.Errorf("team service err: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: %v", model.ErrTeamExists, team.Name)
	}

	team.ID = uuid.New()
	err = s.teamRepo.CreateTeam(ctx, team)
	if err != nil {
		return nil, err
	}

	for _, member := range team.Members {
		userID := uuid.New()
		user := &model.User{
			ID:       userID,
			Username: member.Username,
			TeamID:   team.ID,
			IsActive: member.IsActive,
		}
		err = s.userRepo.CreateUser(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("invalid user_id: %s", user.ID)
		}
	}

	team.Members, err = s.userRepo.GetUsersByTeam(ctx, team.ID)
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (s *TeamService) GetTeamByID(ctx context.Context, id uuid.UUID) (*model.Team, error) {
	team, err := s.teamRepo.GetTeamByID(ctx, id)
	if err != nil {
		return nil, errors.New("team not found")
	}

	// Загружаем участников
	members, err := s.teamRepo.GetMembers(ctx, id)
	if err == nil {
		team.Members = members
	}

	return team, nil
}

func (s *TeamService) UpdateTeam(ctx context.Context, team *model.Team) (*model.Team, error) {
	_, err := s.teamRepo.GetTeamByID(ctx, team.ID)
	if err != nil {
		return nil, errors.New("team not found")
	}

	existing, errGetting := s.teamRepo.GetByName(ctx, team.Name)
	if errGetting != nil {
		return nil, errGetting
	}
	if existing != nil && existing.ID != team.ID {
		return nil, errors.New("team with this name already exists")
	}

	err = s.teamRepo.UpdateTeam(ctx, team)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (s *TeamService) DeleteTeam(ctx context.Context, id uuid.UUID) error {
	return s.teamRepo.DeleteTeam(ctx, id)
}
