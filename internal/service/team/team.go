package team

import (
	"context"
	"errors"

	"github.com/anguless/mr-reviewer/internal/model"
)

func (s *teamService) TeamAddPost(ctx context.Context, team *model.Team) (*model.Team, error) {
	exists, err := s.TeamRepository.TeamExists(ctx, team.TeamName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("team already exists")
	}

	team, err = s.TeamRepository.Create(ctx, team)
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (s *teamService) TeamGetGet(ctx context.Context, teamName string) (*model.Team, error) {
	team, err := s.TeamRepository.GetByName(ctx, teamName)
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (s *teamService) GetTeamNameByUser(ctx context.Context, userID string) (string, error) {
	teamName, err := s.TeamRepository.GetTeamNameByUserID(ctx, userID)
	if err != nil {
		return "", err
	}

	return teamName, nil
}
