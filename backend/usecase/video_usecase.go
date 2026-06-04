package usecase

import (
	"errors"

	"github.com/hinata2398/my-video-app/backend/domain/entity"
	"github.com/hinata2398/my-video-app/backend/domain/repository"
)

type VideoUsecase struct {
	videoRepo repository.VideoRepository
}

func NewVideoUsecase(videoRepo repository.VideoRepository) *VideoUsecase {
	return &VideoUsecase{videoRepo: videoRepo}
}

func (u *VideoUsecase) Create(userID int64, title, description, thumbnailURL string) (*entity.Video, error) {
	if title == "" {
		return nil, errors.New("タイトルは必須です")
	}
	return u.videoRepo.Create(userID, title, description, thumbnailURL)
}

func (u *VideoUsecase) FindAll() ([]*entity.Video, error) {
	return u.videoRepo.FindAll()
}

func (u *VideoUsecase) FindByID(id int64) (*entity.Video, error) {
	video, err := u.videoRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("動画が見つかりません")
	}
	return video, nil
}

func (u *VideoUsecase) Update(id, userID int64, title, description, thumbnailURL string) (*entity.Video, error) {
	if title == "" {
		return nil, errors.New("タイトルは必須です")
	}
	video, err := u.videoRepo.Update(id, userID, title, description, thumbnailURL)
	if err != nil {
		return nil, errors.New("動画が見つからないか、編集権限がありません")
	}
	return video, nil
}

func (u *VideoUsecase) Delete(id, userID int64) error {
	if err := u.videoRepo.Delete(id, userID); err != nil {
		return errors.New("動画が見つからないか、削除権限がありません")
	}
	return nil
}
