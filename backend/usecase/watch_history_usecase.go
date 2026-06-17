package usecase

import (
	"github.com/hinata2398/my-video-app/backend/domain/entity"
	"github.com/hinata2398/my-video-app/backend/domain/repository"
)

type WatchHistoryUsecase struct {
	watchHistoryRepo repository.WatchHistoryRepository
}

func NewWatchHistoryUsecase(watchHistoryRepo repository.WatchHistoryRepository) *WatchHistoryUsecase {
	return &WatchHistoryUsecase{watchHistoryRepo: watchHistoryRepo}
}

func (u *WatchHistoryUsecase) Add(userID, videoID int64) error {
	return u.watchHistoryRepo.Add(userID, videoID)
}

func (u *WatchHistoryUsecase) FindByUserID(userID int64) ([]*entity.Video, error) {
	videos, err := u.watchHistoryRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if videos == nil {
		return []*entity.Video{}, nil
	}
	return videos, nil
}
