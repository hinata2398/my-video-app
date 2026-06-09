package repository

import "github.com/hinata2398/my-video-app/backend/domain/entity"

type UserRepository interface {
	Create(email, passwordHash, username string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindByID(id int64) (*entity.User, error)
	UpdateProfile(id int64, username, avatarURL string) error
}
