package usecase

import (
	"errors"
	"os"
	"testing"

	"github.com/hinata2398/my-video-app/backend/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthUsecase_Register(t *testing.T) {
	t.Run("正常系: 権限を取得できる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		usecase := NewAuthUsecase(mockRepo)
		os.Setenv("JWT_SECRET", "test-secret")

		expected := &entity.User{ID: 1, Email: "test@example.com", Username: "testuser"}
		mockRepo.On("Create", "test@example.com", mock.AnythingOfType("string"), "testuser").Return(expected, nil)

		// Act
		token, err := usecase.Register("test@example.com", "password", "testuser")

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: ユーザー作成がエラーになる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		usecase := NewAuthUsecase(mockRepo)

		mockRepo.On("Create", "", mock.AnythingOfType("string"), "testuser").Return(nil, errors.New("DB error"))

		// Act
		token, err := usecase.Register("", "password", "testuser")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: パスワードが長すぎるとエラーになる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		usecase := NewAuthUsecase(mockRepo)

		// Act
		token, err := usecase.Register("test@example.com", "password123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890", "testuser")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, token)
		mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("異常系: トークン生成エラー", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		usecase := NewAuthUsecase(mockRepo)
		os.Setenv("JWT_SECRET", "")
		t.Cleanup(func() {
			os.Setenv("JWT_SECRET", "test-secret")
		})

		expected := &entity.User{ID: 1, Email: "test@example.com", Username: "testuser"}

		mockRepo.On("Create", "test@example.com", mock.AnythingOfType("string"), "testuser").Return(expected, nil)

		// Act
		token, err := usecase.Register("test@example.com", "password", "testuser")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthUsecase_Login(t *testing.T) {

	t.Run("正常系: ログインできる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		usecase := NewAuthUsecase(mockRepo)
		os.Setenv("JWT_SECRET", "test-secret")

		// bcryptでハッシュ化したパスワードをUserに持たせる
		hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		expected := &entity.User{ID: 1, Email: "test@example.com", PasswordHash: string(hash)}
		mockRepo.On("FindByEmail", "test@example.com").Return(expected, nil)

		// Act
		token, err := usecase.Login("test@example.com", "password")

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: メールアドレスが存在しないとエラーになる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		usecase := NewAuthUsecase(mockRepo)

		mockRepo.On("FindByEmail", "notfound@example.com").Return(nil, errors.New("not found"))

		// Act
		token, err := usecase.Login("notfound@example.com", "password")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.EqualError(t, err, "メールアドレスまたはパスワードが正しくありません")
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: パスワードが間違っているとエラーになる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		usecase := NewAuthUsecase(mockRepo)

		hash, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
		expected := &entity.User{ID: 1, Email: "test@example.com", PasswordHash: string(hash)}
		mockRepo.On("FindByEmail", "test@example.com").Return(expected, nil)

		// Act
		token, err := usecase.Login("test@example.com", "wrong-password")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.EqualError(t, err, "メールアドレスまたはパスワードが正しくありません")
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: トークン生成エラー", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		usecase := NewAuthUsecase(mockRepo)
		os.Setenv("JWT_SECRET", "")
		t.Cleanup(func() {
			os.Setenv("JWT_SECRET", "test-secret")
		})

		hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		expected := &entity.User{ID: 1, Email: "test@example.com", PasswordHash: string(hash)}
		mockRepo.On("FindByEmail", "test@example.com").Return(expected, nil)

		// Act
		token, err := usecase.Login("test@example.com", "password")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
	})
}
