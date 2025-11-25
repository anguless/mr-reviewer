package service

import (
	"context"

	"github.com/anguless/reviewer/internal/model"
	"github.com/anguless/reviewer/internal/repository"
	"github.com/anguless/reviewer/internal/service/pr"
	"github.com/anguless/reviewer/internal/service/stat"
	"github.com/anguless/reviewer/internal/service/team"
	"github.com/anguless/reviewer/internal/service/user"
	"github.com/google/uuid"
)

type Service struct {
	UserService UserService
	TeamService TeamService
	PrService   PrService
	StatService StatService
}

func NewService(repo *repository.Repo) *Service {
	return &Service{
		UserService: user.NewUserService(repo.UserRepo),
		TeamService: team.NewTeamService(repo.TeamRepo, repo.UserRepo),
		PrService:   pr.NewPRService(repo.PrRepo, repo.UserRepo, repo.TeamRepo),
		StatService: stat.NewStatisticsService(repo.StatRepo),
	}
}

type UserService interface {
	CreateUser(ctx context.Context, userC *model.User) (*model.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetAssignedPRs(ctx context.Context, userID uuid.UUID) ([]model.PullRequest, error)
}

type TeamService interface {
	CreateTeam(ctx context.Context, team *model.Team) (*model.Team, error)
	GetTeamByID(ctx context.Context, id uuid.UUID) (*model.Team, error)
	UpdateTeam(ctx context.Context, team *model.Team) (*model.Team, error)
	DeleteTeam(ctx context.Context, id uuid.UUID) error
}

type PrService interface {
	CreatePR(ctx context.Context, pr *model.PullRequest) (*model.PullRequest, error)
	GetPRByID(ctx context.Context, id uuid.UUID) (*model.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID uuid.UUID, oldReviewerID uuid.UUID) (*model.PullRequest, uuid.UUID, error)
	MergePR(ctx context.Context, prID uuid.UUID) (*model.PullRequest, error)
	GetAllPRs(ctx context.Context) ([]model.PullRequest, error)
}

type StatService interface {
	GetReviewStats(ctx context.Context) (*model.ReviewStats, error)
}
