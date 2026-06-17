package usecase

import (
	"github.com/hinata2398/my-video-app/backend/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockWatchHistoryRepository struct {
	mock.Mock
}

func (m *MockWatchHistoryRepository) Add(userID, videoID int64) error {
	args := m.Called(userID, videoID)
	return args.Error(0)
}

func (m *MockWatchHistoryRepository) FindByUserID(userID int64) ([]*entity.Video, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Video), args.Error(1)
}
