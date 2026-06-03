package repository

import "github.com/hinata2398/my-video-app/backend/domain/entity"

type UserRepository interface {
	Create(email, passwordHash string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
}
