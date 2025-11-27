package repository

import (
	"context"

	"github.com/anguless/mr-reviewer/internal/db"
	"github.com/anguless/mr-reviewer/internal/model"
	"github.com/anguless/mr-reviewer/internal/repository/pr"
	"github.com/anguless/mr-reviewer/internal/repository/stat"
	"github.com/anguless/mr-reviewer/internal/repository/team"
	"github.com/anguless/mr-reviewer/internal/repository/user"
	"github.com/google/uuid"
)

type Repo struct {
	UserRepo UserRepository
	TeamRepo TeamRepository
	PrRepo   PrRepository
	StatRepo StatRepository
}

func NewRepository(db db.Database) *Repo {
	return &Repo{
		UserRepo: user.NewUserRepository(db),
		TeamRepo: team.NewTeamRepository(db),
		PrRepo:   pr.NewPRRepository(db),
		StatRepo: stat.NewStatRepository(db),
	}
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetPRsByReviewer(ctx context.Context, userID uuid.UUID) ([]model.PullRequest, error)
	GetActiveUsersByTeam(ctx context.Context, teamID, excludeID uuid.UUID) ([]model.User, error)
	DeactivateTeamMembers(ctx context.Context, teamID uuid.UUID) (int, error)
	GetUsersByTeam(ctx context.Context, teamID uuid.UUID) ([]model.User, error)
	GetActiveUsersByTeamExcluding(ctx context.Context, teamID uuid.UUID, excludeIDs []uuid.UUID) ([]model.User, error)
}

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *model.Team) error
	GetTeamByID(ctx context.Context, id uuid.UUID) (*model.Team, error)
	GetByName(ctx context.Context, name string) (*model.Team, error)
	GetMembers(ctx context.Context, teamID uuid.UUID) ([]model.User, error)
	UpdateTeam(ctx context.Context, team *model.Team) error
	DeleteTeam(ctx context.Context, id uuid.UUID) error
}

type PrRepository interface {
	CreatePR(ctx context.Context, pr *model.PullRequest, reviewers []uuid.UUID) error
	GetPRByID(ctx context.Context, id uuid.UUID) (*model.PullRequest, error)
	UpdatePR(ctx context.Context, pr *model.PullRequest) error
	ReassignReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID uuid.UUID) error
	MergePR(ctx context.Context, prID uuid.UUID) error
	GetAllPRs(ctx context.Context) ([]model.PullRequest, error)
}

type StatRepository interface {
	GetReviewStats(ctx context.Context) (*model.ReviewStats, error)
}
