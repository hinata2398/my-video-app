package usecase

import (
	"github.com/hinata2398/my-video-app/backend/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockCommentRepository struct {
	mock.Mock
}

func (m *MockCommentRepository) Create(videoID, userID int64, body string) (*entity.Comment, error) {
	args := m.Called(videoID, userID, body)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Comment), args.Error(1)
}

func (m *MockCommentRepository) FindByVideoID(videoID, userID int64) ([]*entity.Comment, error) {
	args := m.Called(videoID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Comment), args.Error(1)
}

func (m *MockCommentRepository) ToggleLike(commentID, userID int64) (int64, bool, error) {
	args := m.Called(commentID, userID)
	return args.Get(0).(int64), args.Bool(1), args.Error(2)
}

func (m *MockCommentRepository) ToggleDislike(commentID, userID int64) (int64, bool, error) {
	args := m.Called(commentID, userID)
	return args.Get(0).(int64), args.Bool(1), args.Error(2)
}
