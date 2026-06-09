package handler

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
    "github.com/hinata2398/my-video-app/backend/infrastructure/middleware"
    "github.com/hinata2398/my-video-app/backend/usecase"
)

type LikeHandler struct {
    likeUsecase *usecase.LikeUsecase
}

func NewLikeHandler(likeUsecase *usecase.LikeUsecase) *LikeHandler {
    return &LikeHandler{likeUsecase: likeUsecase}
}

func (h *LikeHandler) Toggle(w http.ResponseWriter, r *http.Request) {
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
    liked, count, err := h.likeUsecase.Toggle(userID, videoID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 結果を返す
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "liked": liked,  // いいねしたか取り消したか
        "count": count,  // 現在のいいね数
    })
}

func (h *LikeHandler) Count(w http.ResponseWriter, r *http.Request) {
    videoID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }

    count, err := h.likeUsecase.Count(videoID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "count": count,
    })
}

func (h *LikeHandler) ToggleDislike(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(middleware.UserIDKey).(int64)
    videoID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }

    disliked, count, err := h.likeUsecase.ToggleDislike(userID, videoID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "disliked": disliked,
        "count":    count,
    })
}

func (h *LikeHandler) DislikeCount(w http.ResponseWriter, r *http.Request) {
    videoID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }

    count, err := h.likeUsecase.DislikeCount(videoID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "count": count,
    })
}

// Status はログイン中ユーザーのいいね・よくないね状態を返す
func (h *LikeHandler) Status(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(middleware.UserIDKey).(int64)
    videoID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }

    liked, _ := h.likeUsecase.LikeExists(userID, videoID)
    disliked, _ := h.likeUsecase.DislikeExists(userID, videoID)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "liked":    liked,
        "disliked": disliked,
    })
}
