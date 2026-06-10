package usecase

import (
	"errors"
	"testing"

	"github.com/hinata2398/my-video-app/backend/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVideoUsecase_Create(t *testing.T) {

	t.Run("正常系: タイトルがあれば動画を作成できる", func(t *testing.T) {
		// Arrange（準備）
		mockRepo := new(MockVideoRepository)
		usecase := NewVideoUsecase(mockRepo)

		expected := &entity.Video{ID: 1, Title: "テスト動画"}
		mockRepo.On("Create", int64(1), "テスト動画", "", "").Return(expected, nil)

		// Act（実行）
		video, err := usecase.Create(1, "テスト動画", "", "")

		// Assert（検証）
		assert.NoError(t, err)
		assert.Equal(t, "テスト動画", video.Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: タイトルが空ならエラーになる", func(t *testing.T) {
		mockRepo := new(MockVideoRepository)
		usecase := NewVideoUsecase(mockRepo)

		video, err := usecase.Create(1, "", "", "")

		assert.Error(t, err)
		assert.Nil(t, video)
		assert.EqualError(t, err, "タイトルは必須です")
		// Repositoryは呼ばれないはず
		mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("異常系: Repositoryがエラーを返したらそのままエラーになる", func(t *testing.T) {
		mockRepo := new(MockVideoRepository)
		usecase := NewVideoUsecase(mockRepo)

		mockRepo.On("Create", int64(1), "テスト動画", "", "").Return(nil, errors.New("DB error"))

		video, err := usecase.Create(1, "テスト動画", "", "")

		assert.Error(t, err)
		assert.Nil(t, video)
		mockRepo.AssertExpectations(t)
	})
}

func TestVideoUsecase_FindByID(t *testing.T) {

	t.Run("正常系: IDで動画を取得できる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockVideoRepository)
		usecase := NewVideoUsecase(mockRepo)

		expected := &entity.Video{ID: 1, Title: "テスト動画"}
		mockRepo.On("FindByID", int64(1)).Return(expected, nil)

		// Act
		video, err := usecase.FindByID(1)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(1), video.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: 存在しないIDならエラーになる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockVideoRepository)
		usecase := NewVideoUsecase(mockRepo)

		mockRepo.On("FindByID", int64(999)).Return(nil, errors.New("not found"))

		// Act
		video, err := usecase.FindByID(999)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, video)
		assert.EqualError(t, err, "動画が見つかりません")
		mockRepo.AssertExpectations(t)
	})
}

func TestVideoUsecase_Update(t *testing.T) {

	t.Run("正常系: タイトルなどを更新できる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockVideoRepository)
		usecase := NewVideoUsecase(mockRepo)

		expected := &entity.Video{
			ID:           1,
			Title:        "テスト動画１",
			Description:  "説明",
			ThumbnailURL: "https://example.com/thumbnail.jpg",
			VideoURL:     "https://example.com/video.mp4",
		}
		mockRepo.On("Update", int64(1), int64(1), "テスト動画１", "説明", "https://example.com/thumbnail.jpg", "https://example.com/video.mp4").Return(expected, nil)

		video, err := usecase.Update(1, 1, "テスト動画１", "説明", "https://example.com/thumbnail.jpg", "https://example.com/video.mp4")

		assert.NoError(t, err)
		assert.Equal(t, int64(1), video.ID)
		assert.Equal(t, "テスト動画１", video.Title)
		assert.Equal(t, "説明", video.Description)
		assert.Equal(t, "https://example.com/thumbnail.jpg", video.ThumbnailURL)
		assert.Equal(t, "https://example.com/video.mp4", video.VideoURL)
	})

	t.Run("異常系: タイトル必須エラー", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockVideoRepository)
		usecase := NewVideoUsecase(mockRepo)

		// Act
		video, err := usecase.Update(1, 1, "", "説明", "https://example.com/thumbnail.jpg", "https://example.com/video.mp4")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, video)
		assert.EqualError(t, err, "タイトルは必須です")
		mockRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("異常系: 動画なし", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockVideoRepository)
		usecase := NewVideoUsecase(mockRepo)

		mockRepo.On("Update", int64(999), int64(1), "テスト動画１", "説明", "https://example.com/thumbnail.jpg", "https://example.com/video.mp4").Return(nil, errors.New("not found"))

		// Act
		video, err := usecase.Update(999, 1, "テスト動画１", "説明", "https://example.com/thumbnail.jpg", "https://example.com/video.mp4")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, video)
		assert.EqualError(t, err, "動画が見つからないか、編集権限がありません")
		mockRepo.AssertExpectations(t)
	})
}

func TestVideoUsecase_Delete(t *testing.T) {

	t.Run("正常系: 動画を削除できる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockVideoRepository)
		usecase := NewVideoUsecase(mockRepo)

		mockRepo.On("Delete", int64(1), int64(1)).Return(nil)

		// Act
		err := usecase.Delete(1, 1)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("異常系: 存在しないIDならエラーになる", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockVideoRepository)
		usecase := NewVideoUsecase(mockRepo)

		mockRepo.On("Delete", int64(999), int64(1)).Return(errors.New("not found"))

		// Act
		err := usecase.Delete(999, 1)

		// Assert
		assert.Error(t, err)
		assert.EqualError(t, err, "動画が見つからないか、削除権限がありません")
		mockRepo.AssertExpectations(t)
	})
}
