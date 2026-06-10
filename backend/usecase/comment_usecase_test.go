package usecase

import (
	"errors"
	"testing"

	"github.com/hinata2398/my-video-app/backend/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCommentUsecase_Create(t *testing.T) {

	t.Run("正常系: コメントを投稿できる", func(t *testing.T) {
		// Arrange（準備）
		mockRepo := new(MockCommentRepository)
		usecase := NewCommentUsecase(mockRepo)

		expected := &entity.Comment{ID: 1, VideoID: 1, UserID: 1, Username: "testuser", Body: "test"}
		mockRepo.On("Create", int64(1), int64(1), "test").Return(expected, nil)

		// Act（実行）
		comment, err := usecase.Create(1, 1, "test")

		// Assert（検証）
		assert.NoError(t, err)
		assert.Equal(t, "test", comment.Body)
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: コメントが空ならエラーになる", func(t *testing.T) {
		mockRepo := new(MockCommentRepository)
		usecase := NewCommentUsecase(mockRepo)

		comment, err := usecase.Create(1, 1, "")

		assert.Error(t, err)
		assert.Nil(t, comment)
		assert.EqualError(t, err, "コメントを入力してください")
		// Repositoryは呼ばれないはず
		mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("異常系: Repositoryがエラーを返したらそのままエラーになる", func(t *testing.T) {
		mockRepo := new(MockCommentRepository)
		usecase := NewCommentUsecase(mockRepo)

		mockRepo.On("Create", int64(999), int64(1), "test").Return(nil, errors.New("DB error"))

		video, err := usecase.Create(999, 1, "test")

		assert.Error(t, err)
		assert.Nil(t, video)
		mockRepo.AssertExpectations(t)
	})
}

func TestCommentUsecase_FindByVideoID(t *testing.T) {

	t.Run("正常系: 動画のコメント一覧を取得できる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockCommentRepository)
		usecase := NewCommentUsecase(mockRepo)

		expected := []*entity.Comment{
			{ID: 1, VideoID: 1, Body: "コメント1"},
			{ID: 2, VideoID: 1, Body: "コメント2"},
		}
		mockRepo.On("FindByVideoID", int64(1), int64(1)).Return(expected, nil)

		// Act
		comments, err := usecase.FindByVideoID(1, 1)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, comments, 2)
		assert.Equal(t, "コメント1", comments[0].Body)
		mockRepo.AssertExpectations(t)
	})

	t.Run("正常系: Repositoryがnilを返したら空配列を返す", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockCommentRepository)
		usecase := NewCommentUsecase(mockRepo)

		mockRepo.On("FindByVideoID", int64(1), int64(1)).Return(nil, nil)

		// Act
		comments, err := usecase.FindByVideoID(1, 1)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, comments)  // nilではなく
		assert.Len(t, comments, 0) // 空配列になっている
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: Repositoryがエラーを返したらエラーになる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockCommentRepository)
		usecase := NewCommentUsecase(mockRepo)

		mockRepo.On("FindByVideoID", int64(1), int64(1)).Return(nil, errors.New("DB error"))

		// Act
		comments, err := usecase.FindByVideoID(1, 1)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, comments)
		mockRepo.AssertExpectations(t)
	})
}
