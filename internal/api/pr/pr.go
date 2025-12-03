package pr

import (
	"context"
	"errors"

	"github.com/anguless/mr-reviewer/internal/converter"
	"github.com/anguless/mr-reviewer/internal/model"
	mr_v1 "github.com/anguless/mr-reviewer/pkg/openapi/mr/v1"
)

func (h *prHandler) PullRequestCreatePost(ctx context.Context, req *mr_v1.PullRequestCreatePostReq) (mr_v1.PullRequestCreatePostRes, error) {
	serviceReq := converter.OpenAPIToServicePRCreateReq(req)
	if serviceReq == nil {
		return &mr_v1.PullRequestCreatePostNotFound{
			Error: mr_v1.ErrorResponseError{
				Code:    mr_v1.ErrorResponseErrorCodeNOTFOUND,
				Message: "Request is nil",
			},
		}, nil
	}

	createdPR, err := h.prService.PullRequestCreatePost(ctx, serviceReq)
	if err != nil {
		return &mr_v1.PullRequestCreatePostConflict{
			Error: mr_v1.ErrorResponseError{
				Code:    mr_v1.ErrorResponseErrorCodePREXISTS,
				Message: err.Error(),
			},
		}, nil
	}

	openAPIPR := converter.ServiceToOpenAPIPR(createdPR)
	if openAPIPR == nil {
		return &mr_v1.PullRequestCreatePostNotFound{
			Error: mr_v1.ErrorResponseError{
				Code:    mr_v1.ErrorResponseErrorCodeNOTFOUND,
				Message: "Created PR is nil",
			},
		}, nil
	}

	optPR := mr_v1.NewOptPullRequest(*openAPIPR)
	result := mr_v1.PullRequestCreatePostCreated{}
	result.SetPr(optPR)

	return &result, nil
}

func (h *prHandler) PullRequestMergePost(ctx context.Context, req *mr_v1.PullRequestMergePostReq) (mr_v1.PullRequestMergePostRes, error) {
	serviceReq := converter.OpenAPIToServicePRMergeReq(req)
	if serviceReq == nil {
		return &mr_v1.ErrorResponse{
			Error: mr_v1.ErrorResponseError{
				Code:    mr_v1.ErrorResponseErrorCodeNOTFOUND,
				Message: "Request is nil",
			},
		}, nil
	}

	mergedPR, err := h.prService.PullRequestMergePost(ctx, serviceReq.PullRequestID)
	if err != nil {
		return &mr_v1.ErrorResponse{
			Error: mr_v1.ErrorResponseError{
				Code:    mr_v1.ErrorResponseErrorCodeNOTFOUND,
				Message: err.Error(),
			},
		}, nil
	}

	openAPIPR := converter.ServiceToOpenAPIPR(mergedPR)
	if openAPIPR == nil {
		return &mr_v1.ErrorResponse{
			Error: mr_v1.ErrorResponseError{
				Code:    mr_v1.ErrorResponseErrorCodeNOTFOUND,
				Message: "Merged PR is nil",
			},
		}, nil
	}

	optPR := mr_v1.NewOptPullRequest(*openAPIPR)
	result := mr_v1.PullRequestMergePostOK{}
	result.SetPr(optPR)

	return &result, nil
}

func (h *prHandler) PullRequestReassignPost(ctx context.Context, req *mr_v1.PullRequestReassignPostReq) (mr_v1.PullRequestReassignPostRes, error) {
	errResp := mr_v1.ErrorResponseError{
		Code:    mr_v1.ErrorResponseErrorCodeNOTFOUND,
		Message: "",
	}
	if req == nil {
		errResp.Message = "Request is nil"
		return &mr_v1.PullRequestReassignPostNotFound{
			Error: errResp,
		}, nil
	}

	reassignedPR, newReviewerID, err := h.prService.PullRequestReassignPost(ctx, req.PullRequestID, req.OldUserID)
	if err != nil {

		errResp.Message = err.Error()

		switch {
		case errors.Is(err, model.ErrMergedPR):
			errResp.Code = mr_v1.ErrorResponseErrorCodePRMERGED
		case errors.Is(err, model.ErrReviewerNotAssigned):
			errResp.Code = mr_v1.ErrorResponseErrorCodeNOTASSIGNED
		case errors.Is(err, model.ErrNoSuitableCandidate):
			errResp.Code = mr_v1.ErrorResponseErrorCodeNOCANDIDATE
		default:
			return &mr_v1.PullRequestReassignPostNotFound{
				Error: errResp,
			}, nil
		}

		return &mr_v1.PullRequestReassignPostConflict{
			Error: errResp,
		}, nil
	}

	openAPIPR := converter.ServiceToOpenAPIPR(reassignedPR)
	if openAPIPR == nil {
		errResp.Message = "Reassigned PR is nil"
		return &mr_v1.PullRequestReassignPostNotFound{
			Error: errResp,
		}, nil
	}

	result := mr_v1.PullRequestReassignPostOK{}
	result.SetPr(*openAPIPR)
	result.SetReplacedBy(newReviewerID)

	return &result, nil
}
