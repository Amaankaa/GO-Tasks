package usecases

import (
	"task-manager/Domain"
)

type UserUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (uu *UserUsecase) RegisterUser(user domain.User) (domain.User, error) {
	return uu.userRepo.RegisterUser(user)
}

func (uu *UserUsecase) LoginUser(user domain.User) (domain.LoginResponse, error) {
	return uu.userRepo.LoginUser(user)
}

func (uu *UserUsecase) PromoteUser(id string) (domain.User, error) {
	return uu.userRepo.PromoteUser(id)
}