package usecase

import (
	"github.com/stretchr/testify/mock"
)

type MockLikeRepository struct {
	mock.Mock
}

func (m *MockLikeRepository) Like(userID, videoID int64) error {
	args := m.Called(userID, videoID)
	return args.Error(0)
}

func (m *MockLikeRepository) Unlike(userID, videoID int64) error {
	args := m.Called(userID, videoID)
	return args.Error(0)
}

func (m *MockLikeRepository) Count(videoID int64) (int64, error) {
	args := m.Called(videoID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockLikeRepository) Exists(userID, videoID int64) (bool, error) {
	args := m.Called(userID, videoID)
	return args.Bool(0), args.Error(1)
}

func (m *MockLikeRepository) Dislike(userID, videoID int64) error {
	args := m.Called(userID, videoID)
	return args.Error(0)
}

func (m *MockLikeRepository) Undislike(userID, videoID int64) error {
	args := m.Called(userID, videoID)
	return args.Error(0)
}

func (m *MockLikeRepository) DislikeCount(videoID int64) (int64, error) {
	args := m.Called(videoID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockLikeRepository) DislikeExists(userID, videoID int64) (bool, error) {
	args := m.Called(userID, videoID)
	return args.Bool(0), args.Error(1)
}
