package pr

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/anguless/mr-reviewer/internal/model"
)

func (r *prRepository) Create(ctx context.Context, pr *model.PullRequest) (*model.PullRequest, error) {
	existingPR, err := r.GetByID(ctx, pr.PullRequestID)
	if err == nil && existingPR != nil {
		return nil, fmt.Errorf("pull request with id %s already exists", pr.PullRequestID)
	}

	query := `
		INSERT INTO pull_requests (id, name, author_id, status, created_at) 
		VALUES ($1, $2, $3, $4, $5)`

	createdAt := pr.CreatedAt
	if createdAt == nil {
		now := time.Now()
		createdAt = &now
	}

	_, err = r.db.Exec(ctx, query, pr.PullRequestID, pr.PullRequestName, pr.AuthorID, pr.Status, createdAt)
	if err != nil {
		return nil, err
	}

	for _, reviewerID := range pr.AssignedReviewers {
		err = r.addReviewer(ctx, pr.PullRequestID, reviewerID)
		if err != nil {
			return nil, err
		}
	}

	return r.GetByID(ctx, pr.PullRequestID)
}

func (r *prRepository) GetByID(ctx context.Context, prID string) (*model.PullRequest, error) {
	query := `SELECT id, name, author_id, status, created_at, merged_at FROM pull_requests WHERE id = $1`
	row := r.db.QueryRow(ctx, query, prID)

	var pr model.PullRequest
	err := row.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.MergedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("pull request with id %s not found", prID)
		}
		return nil, err
	}

	reviewers, err := r.GetReviewers(ctx, prID)
	if err != nil {
		return nil, err
	}
	pr.AssignedReviewers = reviewers

	return &pr, nil
}

func (r *prRepository) UpdateStatus(ctx context.Context, prID string, status model.PrStatus) (*model.PullRequest, error) {
	var query string
	var args []interface{}

	if status == model.PrStatusMerged {
		query = `UPDATE pull_requests SET status = $1, merged_at = $2 WHERE id = $3 RETURNING id, name, author_id, status, created_at, merged_at`
		args = []interface{}{status, time.Now(), prID}
	} else {
		query = `UPDATE pull_requests SET status = $1 WHERE id = $2 RETURNING id, name, author_id, status, created_at, merged_at`
		args = []interface{}{status, prID}
	}

	row := r.db.QueryRow(ctx, query, args...)

	var pr model.PullRequest
	err := row.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.MergedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("pull request with id %s not found", prID)
		}
		return nil, err
	}

	reviewers, err := r.GetReviewers(ctx, prID)
	if err != nil {
		return nil, err
	}
	pr.AssignedReviewers = reviewers

	return &pr, nil
}

func (r *prRepository) UpdateReviewers(ctx context.Context, prID string, reviewers []string) (*model.PullRequest, error) {
	deleteQuery := `DELETE FROM pr_reviewers WHERE pr_id = $1`
	_, err := r.db.Exec(ctx, deleteQuery, prID)
	if err != nil {
		return nil, err
	}

	for _, reviewerID := range reviewers {
		err = r.addReviewer(ctx, prID, reviewerID)
		if err != nil {
			return nil, err
		}
	}

	return r.GetByID(ctx, prID)
}

func (r *prRepository) GetReviewers(ctx context.Context, prID string) ([]string, error) {
	query := `SELECT reviewer_id FROM pr_reviewers WHERE pr_id = $1 ORDER BY reviewer_id`
	rows, err := r.db.Query(ctx, query, prID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var reviewerID string
		err := rows.Scan(&reviewerID)
		if err != nil {
			return nil, err
		}
		reviewers = append(reviewers, reviewerID)
	}

	return reviewers, nil
}

func (r *prRepository) ReplaceReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID string) (*model.PullRequest, error) {
	query := `SELECT 1 FROM pr_reviewers WHERE pr_id = $1 AND reviewer_id = $2 LIMIT 1`
	row := r.db.QueryRow(ctx, query, prID, oldReviewerID)

	var exists int
	err := row.Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("reviewer %s is not assigned to pull request %s", oldReviewerID, prID)
		}
		return nil, err
	}

	deleteQuery := `DELETE FROM pr_reviewers WHERE pr_id = $1 AND reviewer_id = $2`
	_, err = r.db.Exec(ctx, deleteQuery, prID, oldReviewerID)
	if err != nil {
		return nil, err
	}

	err = r.addReviewer(ctx, prID, newReviewerID)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, prID)
}

func (r *prRepository) addReviewer(ctx context.Context, prID, reviewerID string) error {
	query := `INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES ($1, $2)`
	_, err := r.db.Exec(ctx, query, prID, reviewerID)
	return err
}
