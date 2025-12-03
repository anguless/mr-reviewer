package pr

import (
	"context"

	"github.com/anguless/mr-reviewer/internal/db"
	"github.com/anguless/mr-reviewer/internal/model"
)

type prRepository struct {
	db db.Database
}

func NewPRRepository(db db.Database) *prRepository {
	return &prRepository{
		db: db,
	}
}

type PRRepository interface {
	Create(ctx context.Context, pr *model.PullRequest) (*model.PullRequest, error)
	GetByID(ctx context.Context, prID string) (*model.PullRequest, error)
	UpdateStatus(ctx context.Context, prID string, status model.PrStatus) (*model.PullRequest, error)
	UpdateReviewers(ctx context.Context, prID string, reviewers []string) (*model.PullRequest, error)
	GetReviewers(ctx context.Context, prID string) ([]string, error)
	ReplaceReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID string) (*model.PullRequest, error)
}
