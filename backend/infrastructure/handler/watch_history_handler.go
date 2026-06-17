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

type WatchHistoryHandler struct {
	watchHistoryUsecase *usecase.WatchHistoryUsecase
	resolver            MediaURLResolver
}

func NewWatchHistoryHandler(watchHistoryUsecase *usecase.WatchHistoryUsecase, resolver MediaURLResolver) *WatchHistoryHandler {
	return &WatchHistoryHandler{watchHistoryUsecase: watchHistoryUsecase}
}

func (h *WatchHistoryHandler) Add(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)
	videoID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.watchHistoryUsecase.Add(userID, videoID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *WatchHistoryHandler) FindByUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)
	videos, err := h.watchHistoryUsecase.FindByUserID(userID)
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
