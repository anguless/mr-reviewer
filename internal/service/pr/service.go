package pr

import "github.com/anguless/reviewer/internal/repository"

type PrService struct {
	prRepo   repository.PrRepository
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewPRService(prRepo repository.PrRepository, userRepo repository.UserRepository, teamRepo repository.TeamRepository) *PrService {
	return &PrService{
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}
