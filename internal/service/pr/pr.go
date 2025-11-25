package pr

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/anguless/reviewer/internal/model"
	"github.com/google/uuid"
)

func (s *PrService) CreatePR(ctx context.Context, pr *model.PullRequest) (*model.PullRequest, error) {
	author, err := s.userRepo.GetUserByID(ctx, pr.AuthorID)
	if err != nil {
		return nil, errors.New("Автор/команда не найдены")
	}

	candidates, err := s.userRepo.GetActiveUsersByTeam(ctx, author.TeamID, pr.AuthorID)
	if err != nil {
		return nil, err
	}

	reviewers := s.selectRandomReviewers(ctx, candidates, 2)

	pr.ID = uuid.New()
	pr.Status = model.OPEN
	pr.CreatedAt = time.Now()
	pr.Reviewers = reviewers

	err = s.prRepo.CreatePR(ctx, pr, reviewers)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *PrService) GetPRByID(ctx context.Context, id uuid.UUID) (*model.PullRequest, error) {
	pr, err := s.prRepo.GetPRByID(ctx, id)
	if err != nil {
		return nil, errors.New("pull request not found")
	}
	return pr, nil
}

func (s *PrService) ReassignReviewer(
	ctx context.Context,
	prID uuid.UUID,
	oldReviewerID uuid.UUID,
) (*model.PullRequest, uuid.UUID, error) {
	pr, err := s.prRepo.GetPRByID(ctx, prID)
	if err != nil {
		return nil, uuid.Nil, errors.New("pull request not found")
	}

	if pr.Status == model.MERGED {
		return nil, uuid.Nil, errors.New("cannot reassign reviewers for merged PR")
	}

	found := false
	for _, reviewerID := range pr.Reviewers {
		if reviewerID == oldReviewerID {
			found = true
			break
		}
	}
	if !found {
		return nil, uuid.Nil, errors.New("reviewer not assigned to this PR")
	}

	oldReviewer, err := s.userRepo.GetUserByID(ctx, oldReviewerID)
	if err != nil {
		return nil, uuid.Nil, errors.New("old reviewer not found")
	}

	excludeIDs := []uuid.UUID{pr.AuthorID, oldReviewerID}
	for _, reviewerID := range pr.Reviewers {
		if reviewerID != oldReviewerID {
			excludeIDs = append(excludeIDs, reviewerID)
		}
	}

	candidates, err := s.userRepo.GetActiveUsersByTeamExcluding(
		ctx,
		oldReviewer.TeamID,
		excludeIDs,
	)
	if err != nil {
		return nil, uuid.Nil, err
	}
	if len(candidates) == 0 {
		return nil, uuid.Nil, errors.New("no available reviewers in the team")
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	newReviewerID := candidates[rng.Intn(len(candidates))].ID

	err = s.prRepo.ReassignReviewer(ctx, prID, oldReviewerID, newReviewerID)
	if err != nil {
		return nil, uuid.Nil, err
	}

	updated, err := s.prRepo.GetPRByID(ctx, prID)
	if err != nil {
		return nil, uuid.Nil, err
	}

	return updated, newReviewerID, nil
}

func (s *PrService) MergePR(ctx context.Context, prID uuid.UUID) (*model.PullRequest, error) {
	_, err := s.prRepo.GetPRByID(ctx, prID)
	if err != nil {
		return nil, errors.New("pull request not found")
	}

	err = s.prRepo.MergePR(ctx, prID)
	if err != nil {
		return nil, err
	}

	return s.prRepo.GetPRByID(ctx, prID)
}

func (s *PrService) GetAllPRs(ctx context.Context) ([]model.PullRequest, error) {
	return s.prRepo.GetAllPRs(ctx)
}

func (s *PrService) selectRandomReviewers(_ context.Context, candidates []model.User, maxCount int) []uuid.UUID {
	if len(candidates) == 0 {
		return []uuid.UUID{}
	}

	count := maxCount
	if len(candidates) < maxCount {
		count = len(candidates)
	}

	shuffled := make([]model.User, len(candidates))
	copy(shuffled, candidates)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	reviewers := make([]uuid.UUID, 0, count)
	for i := 0; i < count; i++ {
		reviewers = append(reviewers, shuffled[i].ID)
	}

	return reviewers
}
