package user

import (
	"context"

	"github.com/anguless/mr-reviewer/internal/converter"
	mr_v1 "github.com/anguless/mr-reviewer/pkg/openapi/mr/v1"
)

func (h *userHandler) UsersGetReviewGet(ctx context.Context, params mr_v1.UsersGetReviewGetParams) (mr_v1.UsersGetReviewGetRes, error) {
	userID := string(params.UserID)

	prs, err := h.prService.GetPRsByReviewer(ctx, userID)
	if err != nil {
		return &mr_v1.ErrorResponse{
			Error: mr_v1.ErrorResponseError{
				Code:    mr_v1.ErrorResponseErrorCodeNOTFOUND,
				Message: err.Error(),
			},
		}, nil
	}

	openAPIPRList := converter.ServiceToOpenAPIPRList(prs)
	result := mr_v1.UsersGetReviewGetOK{}
	result.SetUserID(userID)
	result.SetPullRequests(openAPIPRList)

	return &result, nil
}

func (h *userHandler) UsersSetIsActivePost(ctx context.Context, req *mr_v1.UsersSetIsActivePostReq) (mr_v1.UsersSetIsActivePostRes, error) {
	serviceUser := converter.OpenAPIToServiceUserSetIsActiveReq(req)
	if serviceUser == nil {
		return &mr_v1.ErrorResponse{
			Error: mr_v1.ErrorResponseError{
				Code:    mr_v1.ErrorResponseErrorCodeNOTFOUND,
				Message: "Request is nil",
			},
		}, nil
	}

	updatedUser, err := h.userService.UsersSetIsActivePost(ctx, serviceUser.ID, serviceUser.IsActive)
	if err != nil {
		return &mr_v1.ErrorResponse{
			Error: mr_v1.ErrorResponseError{
				Code:    mr_v1.ErrorResponseErrorCodeNOTFOUND,
				Message: err.Error(),
			},
		}, nil
	}

	result := converter.ServiceToOpenAPIUserOK(updatedUser)
	if result == nil {
		return &mr_v1.ErrorResponse{
			Error: mr_v1.ErrorResponseError{
				Code:    mr_v1.ErrorResponseErrorCodeNOTFOUND,
				Message: "Updated user is nil",
			},
		}, nil
	}

	return result, nil
}
