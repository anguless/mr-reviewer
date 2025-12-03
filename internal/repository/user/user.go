package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/anguless/mr-reviewer/internal/model"
)

func (r *userRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	query := `
		INSERT INTO users (id, name, is_active, team_name) 
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			is_active = EXCLUDED.is_active,
			team_name = EXCLUDED.team_name
		RETURNING id, name, is_active, team_name`

	var createdUser model.User
	err := r.db.QueryRow(ctx, query, user.ID, user.Name, user.IsActive, user.TeamName).Scan(
		&createdUser.ID, &createdUser.Name, &createdUser.IsActive, &createdUser.TeamName,
	)
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (r *userRepository) GetByID(ctx context.Context, userID string) (*model.User, error) {
	query := `SELECT id, name, is_active, team_name FROM users WHERE id = $1`
	row := r.db.QueryRow(ctx, query, userID)

	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.IsActive, &user.TeamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with id %s not found", userID)
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByTeam(ctx context.Context, teamName string) ([]model.User, error) {
	query := `SELECT id, name, is_active, team_name FROM users WHERE team_name = $1`
	rows, err := r.db.Query(ctx, query, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Name, &user.IsActive, &user.TeamName)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) UpdateIsActive(ctx context.Context, userID string, isActive bool) (*model.User, error) {
	query := `UPDATE users SET is_active = $1 WHERE id = $2 RETURNING id, name, is_active, team_name`
	row := r.db.QueryRow(ctx, query, isActive, userID)

	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.IsActive, &user.TeamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with id %s not found", userID)
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetActiveByTeam(ctx context.Context, teamName string) ([]model.User, error) {
	query := `SELECT id, name, is_active, team_name FROM users WHERE team_name = $1 AND is_active = true`
	rows, err := r.db.Query(ctx, query, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Name, &user.IsActive, &user.TeamName)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) GetAssignedPRs(ctx context.Context, userID string) ([]model.PullRequestShort, error) {
	query := `
		SELECT pr.id, pr.name, pr.author_id, pr.status
		FROM pull_requests pr
		JOIN pr_reviewers prr ON pr.id = prr.pr_id
		WHERE prr.reviewer_id = $1
		ORDER BY pr.created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []model.PullRequestShort
	for rows.Next() {
		var pr model.PullRequestShort
		err := rows.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status)
		if err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}

	return prs, nil
}
