package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hinata2398/my-video-app/backend/domain/entity"
	"github.com/hinata2398/my-video-app/backend/infrastructure/middleware"
	"github.com/hinata2398/my-video-app/backend/usecase"
)

type BookmarkHandler struct {
	bookmarkUsecase *usecase.BookmarkUsecase
	resolver        MediaURLResolver
}

func NewBookmarkHandler(bookmarkUsecase *usecase.BookmarkUsecase, resolver MediaURLResolver) *BookmarkHandler {
	return &BookmarkHandler{bookmarkUsecase: bookmarkUsecase}
}

func (h *BookmarkHandler) Toggle(w http.ResponseWriter, r *http.Request) {
	// ログイン中のユーザーIDをcontextから取り出す
	// （authMiddlewareが事前にセットしてくれている）
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	// URLの {id} 部分（動画ID）を取り出す
	videoID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// usecaseに丸投げ
	bookmarked, err := h.bookmarkUsecase.Toggle(userID, videoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 結果を返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"bookmarked": bookmarked, // いいねしたか取り消したか
	})
}

func (h *BookmarkHandler) Status(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)
	videoID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	bookmarked, err := h.bookmarkUsecase.Exists(userID, videoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"bookmarked": bookmarked,
	})
}

func (h *BookmarkHandler) FindByUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)
	videos, err := h.bookmarkUsecase.FindByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if videos == nil {
		videos = make([]*entity.Video, 0)
	}
	w.Header().Set("Content-Type", "application/json")
	resolveVideos(videos, h.resolver)
	json.NewEncoder(w).Encode(videos)
}
