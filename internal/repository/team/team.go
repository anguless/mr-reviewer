package team

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/anguless/mr-reviewer/internal/model"
)

func (r *teamRepository) Create(ctx context.Context, team *model.Team) (*model.Team, error) {
	query := `
		INSERT INTO users (id, name, is_active, team_name) 
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			is_active = EXCLUDED.is_active,
			team_name = EXCLUDED.team_name
		RETURNING id, name, is_active`

	for _, member := range team.Members {
		var userID, username string
		var isActive bool
		err := r.db.QueryRow(ctx, query, member.UserID, member.Username, member.IsActive, team.TeamName).Scan(
			&userID, &username, &isActive,
		)
		if err != nil {
			return nil, err
		}
	}

	return r.GetByName(ctx, team.TeamName)
}

func (r *teamRepository) GetByName(ctx context.Context, teamName string) (*model.Team, error) {
	query := `SELECT id, name, is_active FROM users WHERE team_name = $1 ORDER BY name`
	rows, err := r.db.Query(ctx, query, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []model.TeamMember
	for rows.Next() {
		var member model.TeamMember
		err := rows.Scan(&member.UserID, &member.Username, &member.IsActive)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	if len(members) == 0 {
		return nil, fmt.Errorf("team with name %s not found", teamName)
	}

	team := &model.Team{
		TeamName: teamName,
		Members:  members,
	}

	return team, nil
}

func (r *teamRepository) TeamExists(ctx context.Context, teamName string) (bool, error) {
	query := `SELECT 1 FROM users WHERE team_name = $1 LIMIT 1`
	row := r.db.QueryRow(ctx, query, teamName)

	var exists int
	err := row.Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *teamRepository) GetTeamNameByUserID(ctx context.Context, userID string) (string, error) {
	query := `SELECT team_name FROM users WHERE id = $1`
	row := r.db.QueryRow(ctx, query, userID)

	var teamName string
	err := row.Scan(&teamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("user with id %s not found", userID)
		}
		return "", err
	}

	return teamName, nil
}
