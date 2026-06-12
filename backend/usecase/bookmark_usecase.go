package usecase

import (
	"github.com/hinata2398/my-video-app/backend/domain/entity"
	"github.com/hinata2398/my-video-app/backend/domain/repository"
)

type BookmarkUsecase struct {
	bookmarkRepo repository.BookmarkRepository
}

func NewBookmarkUsecase(bookmarkRepo repository.BookmarkRepository) *BookmarkUsecase {
	return &BookmarkUsecase{bookmarkRepo: bookmarkRepo}
}

func (u *BookmarkUsecase) Toggle(userID, videoID int64) (bool, error) {
	exists, err := u.bookmarkRepo.Exists(userID, videoID)
	if err != nil {
		return false, err
	}

	if exists {
		err = u.bookmarkRepo.Remove(userID, videoID)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	err = u.bookmarkRepo.Add(userID, videoID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *BookmarkUsecase) Exists(userID, videoID int64) (bool, error) {
	return u.bookmarkRepo.Exists(userID, videoID)
}

func (u *BookmarkUsecase) FindByUserID(userID int64) ([]*entity.Video, error) {
	videos, err := u.bookmarkRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if videos == nil {
		return []*entity.Video{}, nil
	}
	return videos, nil
}
