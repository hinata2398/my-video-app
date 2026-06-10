package usecase

import (
	"errors"
	"testing"

	"github.com/hinata2398/my-video-app/backend/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestUserUsecase_GetMe(t *testing.T) {

	t.Run("正常系: ユーザー情報を取得できる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		expected := &entity.User{ID: 1, Email: "test@example.com", Username: "testuser"}
		mockRepo.On("FindByID", int64(1)).Return(expected, nil)

		// Act
		user, err := usecase.GetMe(1)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(1), user.ID)
		assert.Equal(t, "test@example.com", user.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: 存在しないIDならエラーになる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		mockRepo.On("FindByID", int64(999)).Return(nil, errors.New("not found"))

		// Act
		user, err := usecase.GetMe(999)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_UpdateProfile(t *testing.T) {

	t.Run("正常系: プロフィールを更新できる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		updated := &entity.User{ID: 1, Username: "newname", AvatarURL: "https://example.com/avatar.jpg"}
		mockRepo.On("UpdateProfile", int64(1), "newname", "https://example.com/avatar.jpg").Return(nil)
		mockRepo.On("FindByID", int64(1)).Return(updated, nil)

		// Act
		user, err := usecase.UpdateProfile(1, "newname", "https://example.com/avatar.jpg")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "newname", user.Username)
		assert.Equal(t, "https://example.com/avatar.jpg", user.AvatarURL)
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: UpdateProfileがエラーを返したらエラーになる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		usecase := NewUserUsecase(mockRepo)

		mockRepo.On("UpdateProfile", int64(1), "newname", "").Return(errors.New("DB error"))

		// Act
		user, err := usecase.UpdateProfile(1, "newname", "")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}
