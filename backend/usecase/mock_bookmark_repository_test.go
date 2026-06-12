package usecase

import (
	"github.com/hinata2398/my-video-app/backend/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockBookmarkRepository struct {
	mock.Mock
}

func (m *MockBookmarkRepository) Add(userID, videoID int64) error {
	args := m.Called(userID, videoID)
	return args.Error(0)
}

func (m *MockBookmarkRepository) Remove(userID, videoID int64) error {
	args := m.Called(userID, videoID)
	return args.Error(0)
}

func (m *MockBookmarkRepository) Exists(userID, videoID int64) (bool, error) {
	args := m.Called(userID, videoID)
	return args.Bool(0), args.Error(1)
}

func (m *MockBookmarkRepository) FindByUserID(userID int64) ([]*entity.Video, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Video), args.Error(1)
}
