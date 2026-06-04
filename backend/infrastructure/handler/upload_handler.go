package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hinata2398/my-video-app/backend/infrastructure/middleware"
	"github.com/hinata2398/my-video-app/backend/infrastructure/storage"
)

type UploadHandler struct {
	minio *storage.MinioClient
}

func NewUploadHandler(minio *storage.MinioClient) *UploadHandler {
	return &UploadHandler{minio: minio}
}

func (h *UploadHandler) PresignedURL(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)
	videoID := chi.URLParam(r, "id")

	// オブジェクト名: videos/{userID}/{videoID}/{timestamp}.mp4
	objectName := fmt.Sprintf("videos/%d/%s/%d.mp4", userID, videoID, time.Now().Unix())

	url, err := h.minio.PresignedUploadURL(r.Context(), objectName)
	if err != nil {
		log.Printf("PresignedUploadURL error: %v", err)
		http.Error(w, "failed to generate upload URL", http.StatusInternalServerError)
		return
	}

	publicURL := h.minio.PublicURL(objectName)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"upload_url": url,
		"video_url":  publicURL,
	})
}

func (h *UploadHandler) PresignedThumbnailURL(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)
	objectName := fmt.Sprintf("thumbnails/%d/%d.jpg", userID, time.Now().Unix())

	url, err := h.minio.PresignedUploadURL(r.Context(), objectName)
	if err != nil {
		log.Printf("PresignedUploadURL error: %v", err)
		http.Error(w, "failed to generate upload URL", http.StatusInternalServerError)
		return
	}

	publicURL := h.minio.PublicURL(objectName)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"upload_url":    url,
		"thumbnail_url": publicURL,
	})
}
