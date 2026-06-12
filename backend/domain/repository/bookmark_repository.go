package repository

import "github.com/hinata2398/my-video-app/backend/domain/entity"

type BookmarkRepository interface {
	Add(userID, videoID int64) error
	Remove(userID, videoID int64) error
	Exists(userID, videoID int64) (bool, error)
	FindByUserID(userID int64) ([]*entity.Video, error)
}
