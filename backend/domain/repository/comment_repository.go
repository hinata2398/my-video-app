package repository

import "github.com/hinata2398/my-video-app/backend/domain/entity"

type CommentRepository interface {
	Create(videoID, userID int64, body string) (*entity.Comment, error)
	FindByVideoID(videoID int64) ([]*entity.Comment, error)
}
