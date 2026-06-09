package usecase

import (
	"github.com/hinata2398/my-video-app/backend/domain/entity"
	"github.com/hinata2398/my-video-app/backend/domain/repository"
)

type UserUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

func (u *UserUsecase) GetMe(userID int64) (*entity.User, error) {
	return u.userRepo.FindByID(userID)
}

func (u *UserUsecase) UpdateProfile(userID int64, username, avatarURL string) (*entity.User, error) {
	if err := u.userRepo.UpdateProfile(userID, username, avatarURL); err != nil {
		return nil, err
	}
	return u.userRepo.FindByID(userID)
}
