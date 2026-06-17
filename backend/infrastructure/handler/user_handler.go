package handler

import (
	"encoding/json"
	"net/http"

	"github.com/hinata2398/my-video-app/backend/infrastructure/middleware"
	"github.com/hinata2398/my-video-app/backend/usecase"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
	resolver    MediaURLResolver
}

func NewUserHandler(userUsecase *usecase.UserUsecase, resolver MediaURLResolver) *UserHandler {
	return &UserHandler{userUsecase: userUsecase, resolver: resolver}
}

// GET /api/me - 自分のプロフィール取得
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	user, err := h.userUsecase.GetMe(userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resolveUser(user, h.resolver)
	json.NewEncoder(w).Encode(user)
}

// PUT /api/me - プロフィール更新
func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	var req struct {
		Username  string `json:"username"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.userUsecase.UpdateProfile(userID, req.Username, req.AvatarURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resolveUser(user, h.resolver)
	json.NewEncoder(w).Encode(user)
}
