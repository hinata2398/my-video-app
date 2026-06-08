package usecase

import (
	"errors"

	"github.com/hinata2398/my-video-app/backend/domain/entity"
	"github.com/hinata2398/my-video-app/backend/domain/repository"
)

type CommentUsecase struct {
	commentRepo repository.CommentRepository
}

func NewCommentUsecase(commentRepo repository.CommentRepository) *CommentUsecase {
	return &CommentUsecase{commentRepo: commentRepo}
}

func (u *CommentUsecase) Create(videoID, userID int64, body string) (*entity.Comment, error) {
	if body == "" {
		return nil, errors.New("コメントを入力してください")
	}
	return u.commentRepo.Create(videoID, userID, body)
}

func (u *CommentUsecase) FindByVideoID(videoID int64) ([]*entity.Comment, error) {
	comments, err := u.commentRepo.FindByVideoID(videoID)
	if err != nil {
		return nil, err
	}
	if comments == nil {
		return []*entity.Comment{}, nil
	}
	return comments, nil
}
