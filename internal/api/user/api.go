package user

import (
	"context"

	"github.com/anguless/mr-reviewer/internal/service/pr"
	"github.com/anguless/mr-reviewer/internal/service/user"
	mr_v1 "github.com/anguless/mr-reviewer/pkg/openapi/mr/v1"
)

type userHandler struct {
	userService user.UserService
	prService   pr.PRService
}

func NewUserHandler(userService user.UserService, prService pr.PRService) *userHandler {
	return &userHandler{
		userService: userService,
		prService:   prService,
	}
}

type UserHandler interface {
	UsersGetReviewGet(ctx context.Context, params mr_v1.UsersGetReviewGetParams) (mr_v1.UsersGetReviewGetRes, error)
	UsersSetIsActivePost(ctx context.Context, req *mr_v1.UsersSetIsActivePostReq) (mr_v1.UsersSetIsActivePostRes, error)
}
