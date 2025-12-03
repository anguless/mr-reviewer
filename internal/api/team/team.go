package team

import (
	"context"

	"github.com/anguless/mr-reviewer/internal/converter"
	mr_v1 "github.com/anguless/mr-reviewer/pkg/openapi/mr/v1"
)

func (h *teamHandler) TeamAddPost(ctx context.Context, req *mr_v1.Team) (mr_v1.TeamAddPostRes, error) {
	serviceTeam := converter.OpenAPIToServiceTeam(req)
	if serviceTeam == nil {
		return &mr_v1.ErrorResponse{
			Error: mr_v1.ErrorResponseError{
				Code:    mr_v1.ErrorResponseErrorCodeNOTFOUND,
				Message: "Request is nil",
			},
		}, nil
	}

	createdTeam, err := h.teamService.TeamAddPost(ctx, serviceTeam)
	if err != nil {
		return &mr_v1.ErrorResponse{
			Error: mr_v1.ErrorResponseError{
				Code:    mr_v1.ErrorResponseErrorCodeTEAMEXISTS,
				Message: err.Error(),
			},
		}, nil
	}

	openAPITeam := converter.ServiceToOpenAPITeam(createdTeam)
	if openAPITeam == nil {
		return &mr_v1.ErrorResponse{
			Error: mr_v1.ErrorResponseError{
				Code:    mr_v1.ErrorResponseErrorCodeNOTFOUND,
				Message: "Created team is nil",
			},
		}, nil
	}

	optTeam := mr_v1.NewOptTeam(*openAPITeam)
	result := mr_v1.TeamAddPostCreated{}
	result.SetTeam(optTeam)

	return &result, nil
}

func (h *teamHandler) TeamGetGet(ctx context.Context, params mr_v1.TeamGetGetParams) (mr_v1.TeamGetGetRes, error) {
	teamName := string(params.TeamName)

	team, err := h.teamService.TeamGetGet(ctx, teamName)
	if err != nil {
		return &mr_v1.ErrorResponse{
			Error: mr_v1.ErrorResponseError{
				Code:    mr_v1.ErrorResponseErrorCodeNOTFOUND,
				Message: err.Error(),
			},
		}, nil
	}

	openAPITeam := converter.ServiceToOpenAPITeam(team)
	if openAPITeam == nil {
		return &mr_v1.ErrorResponse{
			Error: mr_v1.ErrorResponseError{
				Code:    mr_v1.ErrorResponseErrorCodeNOTFOUND,
				Message: "Team is nil",
			},
		}, nil
	}

	return openAPITeam, nil
}
