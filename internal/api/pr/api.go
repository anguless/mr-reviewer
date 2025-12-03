package pr

import (
	"context"

	"github.com/anguless/mr-reviewer/internal/service/pr"
	mr_v1 "github.com/anguless/mr-reviewer/pkg/openapi/mr/v1"
)

type prHandler struct {
	prService pr.PRService
}

func NewPRHandler(prService pr.PRService) *prHandler {
	return &prHandler{
		prService: prService,
	}
}

type PRHandler interface {
	PullRequestCreatePost(ctx context.Context, req *mr_v1.PullRequestCreatePostReq) (mr_v1.PullRequestCreatePostRes, error)
	PullRequestMergePost(ctx context.Context, req *mr_v1.PullRequestMergePostReq) (mr_v1.PullRequestMergePostRes, error)
	PullRequestReassignPost(ctx context.Context, req *mr_v1.PullRequestReassignPostReq) (mr_v1.PullRequestReassignPostRes, error)
}
