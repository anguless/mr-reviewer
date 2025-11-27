package user

import "github.com/anguless/mr-reviewer/internal/repository"

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{userRepo: repo}
}
