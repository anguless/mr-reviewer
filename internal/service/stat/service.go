package stat

import "github.com/anguless/mr-reviewer/internal/repository"

type StatService struct {
	statRepo repository.StatRepository
}

func NewStatisticsService(statsRepo repository.StatRepository) *StatService {
	return &StatService{statRepo: statsRepo}
}
