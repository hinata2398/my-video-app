package usecase

import (
	"errors"
	"testing"

	"github.com/hinata2398/my-video-app/backend/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkUsecase_Toggle(t *testing.T) {

	t.Run("正常系: 未ブックマークならブックマークを追加できる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockBookmarkRepository)
		usecase := NewBookmarkUsecase(mockRepo)

		mockRepo.On("Exists", int64(1), int64(1)).Return(false, nil)
		mockRepo.On("Add", int64(1), int64(1)).Return(nil)

		// Act
		bookmarked, err := usecase.Toggle(1, 1)

		// Assert
		assert.NoError(t, err)
		assert.True(t, bookmarked)
		mockRepo.AssertExpectations(t)
	})

	t.Run("正常系: ブックマーク済みならブックマークを解除できる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockBookmarkRepository)
		usecase := NewBookmarkUsecase(mockRepo)

		mockRepo.On("Exists", int64(1), int64(1)).Return(true, nil)
		mockRepo.On("Remove", int64(1), int64(1)).Return(nil)

		// Act
		bookmarked, err := usecase.Toggle(1, 1)

		// Assert
		assert.NoError(t, err)
		assert.False(t, bookmarked)
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: Existsがエラーを返したらエラーになる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockBookmarkRepository)
		usecase := NewBookmarkUsecase(mockRepo)

		mockRepo.On("Exists", int64(1), int64(1)).Return(false, errors.New("DB error"))

		// Act
		bookmarked, err := usecase.Toggle(1, 1)

		// Assert
		assert.Error(t, err)
		assert.False(t, bookmarked)
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: Addがエラーを返したらエラーになる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockBookmarkRepository)
		usecase := NewBookmarkUsecase(mockRepo)

		mockRepo.On("Exists", int64(1), int64(1)).Return(false, nil)
		mockRepo.On("Add", int64(1), int64(1)).Return(errors.New("DB error"))

		// Act
		bookmarked, err := usecase.Toggle(1, 1)

		// Assert
		assert.Error(t, err)
		assert.False(t, bookmarked)
		mockRepo.AssertExpectations(t)
	})
}

func TestBookmarkUsecase_FindByUserID(t *testing.T) {

	t.Run("正常系: ブックマーク一覧を取得できる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockBookmarkRepository)
		usecase := NewBookmarkUsecase(mockRepo)

		expected := []*entity.Video{
			{ID: 1, Title: "動画1"},
			{ID: 2, Title: "動画2"},
		}
		mockRepo.On("FindByUserID", int64(1)).Return(expected, nil)

		// Act
		videos, err := usecase.FindByUserID(1)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, videos, 2)
		assert.Equal(t, "動画1", videos[0].Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("正常系: Repositoryがnilを返したら空配列を返す", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockBookmarkRepository)
		usecase := NewBookmarkUsecase(mockRepo)

		mockRepo.On("FindByUserID", int64(1)).Return(nil, nil)

		// Act
		videos, err := usecase.FindByUserID(1)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, videos)
		assert.Len(t, videos, 0)
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: Repositoryがエラーを返したらエラーになる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockBookmarkRepository)
		usecase := NewBookmarkUsecase(mockRepo)

		mockRepo.On("FindByUserID", int64(1)).Return(nil, errors.New("DB error"))

		// Act
		videos, err := usecase.FindByUserID(1)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, videos)
		mockRepo.AssertExpectations(t)
	})
}
