package repository

import "github.com/hinata2398/my-video-app/backend/domain/entity"

type WatchHistoryRepository interface {
	Add(userID, videoID int64) error
	FindByUserID(userID int64) ([]*entity.Video, error)
}
