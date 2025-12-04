package pr

import (
	"context"

	"github.com/anguless/mr-reviewer/internal/model"
)

func (s *prService) PullRequestCreatePost(ctx context.Context, pr *model.PullRequest) (*model.PullRequest, error) {
	pr.Status = model.PrStatusOpen

	authorTeamName, err := s.TeamRepository.GetTeamNameByUserID(ctx, pr.AuthorID)
	if err != nil {
		return nil, err
	}

	activeUsers, err := s.UserRepository.GetActiveByTeam(ctx, authorTeamName)
	if err != nil {
		return nil, err
	}

	var candidates []model.User
	for _, user := range activeUsers {
		if user.ID != pr.AuthorID {
			candidates = append(candidates, user)
		}
	}

	var reviewers []string
	for i := 0; i < 2 && i < len(candidates); i++ {
		reviewers = append(reviewers, candidates[i].ID)
	}

	pr.AssignedReviewers = reviewers

	createdPR, err := s.PRRepository.Create(ctx, pr)
	if err != nil {
		return nil, err
	}

	return createdPR, nil
}

func (s *prService) PullRequestMergePost(ctx context.Context, prID string) (*model.PullRequest, error) {
	existingPR, err := s.PRRepository.GetByID(ctx, prID)
	if err != nil {
		return nil, err
	}

	if existingPR.Status == model.PrStatusMerged {
		return existingPR, nil
	}

	updatedPR, err := s.PRRepository.UpdateStatus(ctx, prID, model.PrStatusMerged)
	if err != nil {
		return nil, err
	}

	return updatedPR, nil
}

func (s *prService) PullRequestReassignPost(ctx context.Context, prID string, oldUserID string) (*model.PullRequest, string, error) {
	var newCandidate string

	pr, err := s.PRRepository.GetByID(ctx, prID)
	if err != nil {
		return nil, newCandidate, err
	}

	if pr.Status == model.PrStatusMerged {
		return nil, newCandidate, model.ErrMergedPR
	}

	isReviewer := false
	for _, reviewerID := range pr.AssignedReviewers {
		if reviewerID == oldUserID {
			isReviewer = true
			break
		}
	}

	if !isReviewer {
		return nil, newCandidate, model.ErrReviewerNotAssigned
	}

	reviewerTeamName, err := s.TeamRepository.GetTeamNameByUserID(ctx, oldUserID)
	if err != nil {
		return nil, newCandidate, err
	}

	activeUsers, err := s.UserRepository.GetActiveByTeam(ctx, reviewerTeamName)
	if err != nil {
		return nil, newCandidate, err
	}

	var candidates []model.User
	for _, user := range activeUsers {
		isAssigned := false
		for _, assignedReviewer := range pr.AssignedReviewers {
			if user.ID == assignedReviewer || user.ID == pr.AuthorID {
				isAssigned = true
				break
			}
		}
		if !isAssigned && user.ID != oldUserID {
			candidates = append(candidates, user)
		}
	}

	if len(candidates) == 0 {
		return nil, newCandidate, model.ErrNoSuitableCandidate
	}

	newReviewer := candidates[s.rnd.Intn(len(candidates))]

	updatedReviewers := make([]string, len(pr.AssignedReviewers))
	copy(updatedReviewers, pr.AssignedReviewers)
	for i, reviewerID := range updatedReviewers {
		if reviewerID == oldUserID {
			updatedReviewers[i] = newReviewer.ID
			break
		}
	}

	newCandidate = newReviewer.ID

	pr.AssignedReviewers = updatedReviewers
	updatedPR, err := s.PRRepository.UpdateReviewers(ctx, prID, updatedReviewers)
	if err != nil {
		return nil, newCandidate, err
	}

	return updatedPR, newCandidate, nil
}

func (s *prService) GetPRsByReviewer(ctx context.Context, userID string) ([]model.PullRequestShort, error) {
	prs, err := s.UserRepository.GetAssignedPRs(ctx, userID)
	if err != nil {
		return nil, err
	}

	return prs, nil
}
