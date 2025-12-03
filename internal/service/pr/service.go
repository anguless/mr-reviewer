package pr

import (
	"context"

	"github.com/anguless/mr-reviewer/internal/model"
	"github.com/anguless/mr-reviewer/internal/repository/pr"
	"github.com/anguless/mr-reviewer/internal/repository/team"
	"github.com/anguless/mr-reviewer/internal/repository/user"
)

type prService struct {
	PRRepository   pr.PRRepository
	UserRepository user.UserRepository
	TeamRepository team.TeamRepository
}

func NewPRService(prRepo pr.PRRepository, userRepo user.UserRepository, teamRepo team.TeamRepository) *prService {
	return &prService{
		PRRepository:   prRepo,
		UserRepository: userRepo,
		TeamRepository: teamRepo,
	}
}

type PRService interface {
	PullRequestCreatePost(ctx context.Context, pr *model.PullRequest) (*model.PullRequest, error)
	PullRequestMergePost(ctx context.Context, prID string) (*model.PullRequest, error)
	PullRequestReassignPost(ctx context.Context, prID string, oldUserID string) (*model.PullRequest, string, error)
	GetPRsByReviewer(ctx context.Context, userID string) ([]model.PullRequestShort, error)
}
