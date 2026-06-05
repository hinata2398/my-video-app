package usecase

import (
	"github.com/hinata2398/my-video-app/backend/domain/repository"
)

type LikeUsecase struct {
	likeRepo repository.LikeRepository
}

func NewLikeUsecase(likeRepo repository.LikeRepository) *LikeUsecase {
	return &LikeUsecase{likeRepo: likeRepo}
}

func (u *LikeUsecase) Toggle(userID, videoID int64) (liked bool, count int64, err error) {
    exists, err := u.likeRepo.Exists(userID, videoID)
    if err != nil {
        return false, 0, err
    }

    if exists {
        err = u.likeRepo.Unlike(userID, videoID)
        liked = false
    } else {
        err = u.likeRepo.Like(userID, videoID)
        liked = true
    }
    if err != nil {
        return false, 0, err
    }

    count, err = u.likeRepo.Count(videoID)
    return liked, count, err
}

func (u *LikeUsecase) Count(videoID int64) (int64, error) {
    return u.likeRepo.Count(videoID)
}
