package converter

import (
	"time"

	"github.com/anguless/mr-reviewer/internal/model"
	mr_v1 "github.com/anguless/mr-reviewer/pkg/openapi/mr/v1"
)

func OpenAPIToServicePR(pr *mr_v1.PullRequest) *model.PullRequest {
	if pr == nil {
		return nil
	}

	var createdAt *time.Time
	if val, ok := pr.CreatedAt.Get(); ok {
		createdAt = &val
	}

	var mergedAt *time.Time
	if val, ok := pr.MergedAt.Get(); ok {
		mergedAt = &val
	}

	status := model.PrStatus(pr.Status)

	return &model.PullRequest{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            status,
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         createdAt,
		MergedAt:          mergedAt,
	}
}

func OpenAPIToServicePRCreateReq(req *mr_v1.PullRequestCreatePostReq) *model.PullRequest {
	if req == nil {
		return nil
	}

	return &model.PullRequest{
		PullRequestID:   req.PullRequestID,
		PullRequestName: req.PullRequestName,
		AuthorID:        req.AuthorID,
		Status:          model.PrStatusOpen, // Default status for new PR
	}
}

func OpenAPIToServicePRMergeReq(req *mr_v1.PullRequestMergePostReq) *model.PullRequest {
	if req == nil {
		return nil
	}

	return &model.PullRequest{
		PullRequestID: req.PullRequestID,
	}
}

func OpenAPIToServicePRReassignReq(req *mr_v1.PullRequestReassignPostReq) *model.PullRequest {
	if req == nil {
		return nil
	}

	return &model.PullRequest{
		PullRequestID: req.PullRequestID,
	}
}

func OpenAPIToServiceTeam(team *mr_v1.Team) *model.Team {
	if team == nil {
		return nil
	}

	members := make([]model.TeamMember, len(team.Members))
	for i, member := range team.Members {
		members[i] = model.TeamMember{
			UserID:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		}
	}

	return &model.Team{
		TeamName: team.TeamName,
		Members:  members,
	}
}

func OpenAPIToServiceUser(user *mr_v1.User) *model.User {
	if user == nil {
		return nil
	}

	return &model.User{
		ID:       user.UserID,
		Name:     user.Username,
		IsActive: user.IsActive,
		TeamName: user.TeamName,
	}
}

func OpenAPIToServiceUserSetIsActiveReq(req *mr_v1.UsersSetIsActivePostReq) *model.User {
	if req == nil {
		return nil
	}

	return &model.User{
		ID:       req.UserID,
		IsActive: req.IsActive,
	}
}

func ServiceToOpenAPIPR(pr *model.PullRequest) *mr_v1.PullRequest {
	if pr == nil {
		return nil
	}

	var createdAt mr_v1.OptNilDateTime
	if pr.CreatedAt != nil {
		createdAt.SetTo(*pr.CreatedAt)
	}

	var mergedAt mr_v1.OptNilDateTime
	if pr.MergedAt != nil {
		mergedAt.SetTo(*pr.MergedAt)
	}

	status := mr_v1.PullRequestStatus(pr.Status)

	return &mr_v1.PullRequest{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            status,
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         createdAt,
		MergedAt:          mergedAt,
	}
}

func ServiceToOpenAPIPRShort(pr *model.PullRequestShort) *mr_v1.PullRequestShort {
	if pr == nil {
		return nil
	}

	status := mr_v1.PullRequestShortStatus(pr.Status)

	return &mr_v1.PullRequestShort{
		PullRequestID:   pr.PullRequestID,
		PullRequestName: pr.PullRequestName,
		AuthorID:        pr.AuthorID,
		Status:          status,
	}
}

func ServiceToOpenAPIPRList(prList []model.PullRequestShort) []mr_v1.PullRequestShort {
	result := make([]mr_v1.PullRequestShort, len(prList))
	for i, pr := range prList {
		result[i] = *ServiceToOpenAPIPRShort(&pr)
	}
	return result
}

func ServiceToOpenAPITeam(team *model.Team) *mr_v1.Team {
	if team == nil {
		return nil
	}

	members := make([]mr_v1.TeamMember, len(team.Members))
	for i, member := range team.Members {
		members[i] = mr_v1.TeamMember{
			UserID:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		}
	}

	return &mr_v1.Team{
		TeamName: team.TeamName,
		Members:  members,
	}
}

func ServiceToOpenAPIUser(user *model.User) *mr_v1.User {
	if user == nil {
		return nil
	}

	return &mr_v1.User{
		UserID:   user.ID,
		Username: user.Name,
		IsActive: user.IsActive,
		TeamName: user.TeamName,
	}
}

func ServiceToOpenAPIUserOK(user *model.User) *mr_v1.UsersSetIsActivePostOK {
	if user == nil {
		return nil
	}

	optUser := mr_v1.NewOptUser(*ServiceToOpenAPIUser(user))
	result := mr_v1.UsersSetIsActivePostOK{}
	result.SetUser(optUser)
	return &result
}
