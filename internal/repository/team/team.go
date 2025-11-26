package team

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/anguless/reviewer/internal/model"
)

func (r *teamRepository) CreateTeam(ctx context.Context, team *model.Team) error {
	query := `
		INSERT INTO teams (id, name, created_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(ctx, query, team.ID, team.Name, time.Now())
	return err
}

func (r *teamRepository) GetTeamByID(ctx context.Context, id uuid.UUID) (*model.Team, error) {
	query := `
		SELECT id, name
		FROM teams
		WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, id)
	var team model.Team
	err := row.Scan(&team.ID, &team.Name)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *teamRepository) GetByName(ctx context.Context, name string) (*model.Team, error) {
	query := `
		SELECT id, name
		FROM teams
		WHERE name = $1
	`
	row := r.db.QueryRow(ctx, query, name)
	var team model.Team
	err := row.Scan(&team.ID, &team.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &team, nil
}

func (r *teamRepository) GetMembers(ctx context.Context, teamID uuid.UUID) ([]model.User, error) {
	query := `
		SELECT id, username, team_id, is_active
		FROM users
		WHERE team_id = $1
		ORDER BY username
	`
	rows, err := r.db.Query(ctx, query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.ID, &u.Username, &u.TeamID, &u.IsActive)
		if err != nil {
			return nil, err
		}
		members = append(members, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return members, nil
}

func (r *teamRepository) UpdateTeam(ctx context.Context, team *model.Team) error {
	query := `
		UPDATE teams
		SET name = $1
		WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, team.Name, team.ID)
	return err
}

func (r *teamRepository) DeleteTeam(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM teams WHERE id = $1", id)
	return err
}
