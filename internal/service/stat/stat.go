package stat

import (
	"context"

	"github.com/anguless/mr-reviewer/internal/model"
)

func (s *StatService) GetReviewStats(ctx context.Context) (*model.ReviewStats, error) {
	return s.statRepo.GetReviewStats(ctx)
}
