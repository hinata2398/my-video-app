package repository

import "github.com/hinata2398/my-video-app/backend/domain/entity"

type CommentRepository interface {
	Create(videoID, userID int64, body string) (*entity.Comment, error)
	FindByVideoID(videoID, userID int64) ([]*entity.Comment, error)
	ToggleLike(commentID, userID int64) (int64, bool, error)
	ToggleDislike(commentID, userID int64) (int64, bool, error)
}
