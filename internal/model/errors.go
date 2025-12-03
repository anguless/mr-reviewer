package model

import "errors"

var (
	ErrMergedPR            = errors.New("cannot reassign on merged PR")
	ErrReviewerNotAssigned = errors.New("reviewer is not assigned to this PR")
	ErrNoSuitableCandidate = errors.New("no active replacement candidate in team")
)
