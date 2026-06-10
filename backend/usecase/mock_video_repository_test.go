package usecase

import (
	"github.com/hinata2398/my-video-app/backend/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockVideoRepository struct {
	mock.Mock
}

func (m *MockVideoRepository) Create(userID int64, title, description, thumbnailURL string) (*entity.Video, error) {
	args := m.Called(userID, title, description, thumbnailURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Video), args.Error(1)
}

func (m *MockVideoRepository) FindAll(query string) ([]*entity.Video, error) {
	args := m.Called(query)
	return args.Get(0).([]*entity.Video), args.Error(1)
}

func (m *MockVideoRepository) FindByUserID(userID int64) ([]*entity.Video, error) {
	args := m.Called(userID)
	return args.Get(0).([]*entity.Video), args.Error(1)
}

func (m *MockVideoRepository) FindByID(id int64) (*entity.Video, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Video), args.Error(1)
}

func (m *MockVideoRepository) Update(id, userID int64, title, description, thumbnailURL, videoURL string) (*entity.Video, error) {
	args := m.Called(id, userID, title, description, thumbnailURL, videoURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Video), args.Error(1)
}

func (m *MockVideoRepository) UpdateVideoURL(id int64, videoURL string) error {
	args := m.Called(id, videoURL)
	return args.Error(0)
}

func (m *MockVideoRepository) Delete(id, userID int64) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *MockVideoRepository) IncrementViewCount(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
