package repository

import "github.com/hinata2398/my-video-app/backend/domain/entity"

type VideoRepository interface {
	Create(userID int64, title, description, thumbnailURL string) (*entity.Video, error)
	FindAll() ([]*entity.Video, error)
	FindByID(id int64) (*entity.Video, error)
	Update(id, userID int64, title, description, thumbnailURL string) (*entity.Video, error)
	Delete(id, userID int64) error
}
