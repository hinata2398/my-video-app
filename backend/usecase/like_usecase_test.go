package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLikeUsecase_Toggle(t *testing.T) {

	t.Run("正常系: いいね済みならいいねを解除できる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockLikeRepository)
		usecase := NewLikeUsecase(mockRepo)

		mockRepo.On("Exists", int64(1), int64(1)).Return(true, nil)
		mockRepo.On("Unlike", int64(1), int64(1)).Return(nil)
		mockRepo.On("Count", int64(1)).Return(int64(0), nil)

		// Act
		liked, count, err := usecase.Toggle(1, 1)

		// Assert
		assert.NoError(t, err)
		assert.False(t, liked)
		assert.Equal(t, int64(0), count)
		mockRepo.AssertExpectations(t)
	})

	t.Run("正常系: 未いいねならいいねできる（dislikeも解除される）", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockLikeRepository)
		usecase := NewLikeUsecase(mockRepo)

		mockRepo.On("Exists", int64(1), int64(1)).Return(false, nil)
		mockRepo.On("Undislike", int64(1), int64(1)).Return(nil)
		mockRepo.On("Like", int64(1), int64(1)).Return(nil)
		mockRepo.On("Count", int64(1)).Return(int64(1), nil)

		// Act
		liked, count, err := usecase.Toggle(1, 1)

		// Assert
		assert.NoError(t, err)
		assert.True(t, liked)
		assert.Equal(t, int64(1), count)
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: Existsがエラーを返したらエラーになる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockLikeRepository)
		usecase := NewLikeUsecase(mockRepo)

		mockRepo.On("Exists", int64(1), int64(1)).Return(false, errors.New("DB error"))

		// Act
		liked, count, err := usecase.Toggle(1, 1)

		// Assert
		assert.Error(t, err)
		assert.False(t, liked)
		assert.Equal(t, int64(0), count)
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: Unlikeがエラーを返したらエラーになる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockLikeRepository)
		usecase := NewLikeUsecase(mockRepo)

		mockRepo.On("Exists", int64(1), int64(1)).Return(true, nil)
		mockRepo.On("Unlike", int64(1), int64(1)).Return(errors.New("DB error"))

		// Act
		liked, count, err := usecase.Toggle(1, 1)

		// Assert
		assert.Error(t, err)
		assert.False(t, liked)
		assert.Equal(t, int64(0), count)
		mockRepo.AssertExpectations(t)
	})
}

func TestLikeUsecase_ToggleDislike(t *testing.T) {

	t.Run("正常系: dislike済みならdislikeを解除できる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockLikeRepository)
		usecase := NewLikeUsecase(mockRepo)

		mockRepo.On("DislikeExists", int64(1), int64(1)).Return(true, nil)
		mockRepo.On("Undislike", int64(1), int64(1)).Return(nil)
		mockRepo.On("DislikeCount", int64(1)).Return(int64(0), nil)

		// Act
		disliked, count, err := usecase.ToggleDislike(1, 1)

		// Assert
		assert.NoError(t, err)
		assert.False(t, disliked)
		assert.Equal(t, int64(0), count)
		mockRepo.AssertExpectations(t)
	})

	t.Run("正常系: 未dislikeならdislikeできる（likeも解除される）", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockLikeRepository)
		usecase := NewLikeUsecase(mockRepo)

		mockRepo.On("DislikeExists", int64(1), int64(1)).Return(false, nil)
		mockRepo.On("Unlike", int64(1), int64(1)).Return(nil)
		mockRepo.On("Dislike", int64(1), int64(1)).Return(nil)
		mockRepo.On("DislikeCount", int64(1)).Return(int64(1), nil)

		// Act
		disliked, count, err := usecase.ToggleDislike(1, 1)

		// Assert
		assert.NoError(t, err)
		assert.True(t, disliked)
		assert.Equal(t, int64(1), count)
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: DislikeExistsがエラーを返したらエラーになる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockLikeRepository)
		usecase := NewLikeUsecase(mockRepo)

		mockRepo.On("DislikeExists", int64(1), int64(1)).Return(false, errors.New("DB error"))

		// Act
		disliked, count, err := usecase.ToggleDislike(1, 1)

		// Assert
		assert.Error(t, err)
		assert.False(t, disliked)
		assert.Equal(t, int64(0), count)
		mockRepo.AssertExpectations(t)
	})
}
