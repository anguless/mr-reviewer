package team

import (
	"context"
	"errors"
	"fmt"

	"github.com/anguless/reviewer/internal/model"
	"github.com/anguless/reviewer/internal/service/pr"
	"github.com/google/uuid"
)

func (s *TeamService) CreateTeam(ctx context.Context, team *model.Team) (*model.Team, error) {
	existing, err := s.teamRepo.GetByName(ctx, team.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("team with this name already exists" + team.Name)
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

func (s *TeamService) DeactivateTeamMembers(ctx context.Context, teamID uuid.UUID, PrService *pr.PrService) (int, error) {
	_, err := s.teamRepo.GetTeamByID(ctx, teamID)
	if err != nil {
		return 0, errors.New("team not found")
	}

	activeUsers, err := s.userRepo.GetActiveUsersByTeam(ctx, teamID, uuid.Nil)
	if err != nil {
		return 0, err
	}

	if len(activeUsers) == 0 {
		return 0, nil
	}

	allPRs, err := PrService.GetAllPRs(ctx)
	if err != nil {
		return 0, err
	}

	userIDMap := make(map[uuid.UUID]bool)
	for _, user := range activeUsers {
		userIDMap[user.ID] = true
	}

	for _, pr := range allPRs {
		if pr.Status != model.OPEN {
			continue
		}

		for _, reviewerID := range pr.Reviewers {
			if userIDMap[reviewerID] {
				_, _, err = PrService.ReassignReviewer(ctx, pr.ID, reviewerID)
				if err != nil {
					continue
				}
			}
		}
	}

	deactivatedCount, err := s.userRepo.DeactivateTeamMembers(ctx, teamID)
	if err != nil {
		return 0, err
	}

	return deactivatedCount, nil
}
