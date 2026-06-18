package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hinata2398/my-video-app/backend/infrastructure/middleware"
	"github.com/hinata2398/my-video-app/backend/infrastructure/storage"
)

type ThumbnailHandler struct {
	minio    *storage.MinioClient
	db       *sql.DB
	resolver MediaURLResolver
}

func NewThumbnailHandler(minio *storage.MinioClient, db *sql.DB, resolver MediaURLResolver) *ThumbnailHandler {
	return &ThumbnailHandler{minio: minio, db: db, resolver: resolver}
}

func (h *ThumbnailHandler) Generate(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)
	videoID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// 動画URLをDBから取得（自分の動画のみ）
	var videoURL string
	err = h.db.QueryRowContext(r.Context(),
		`SELECT video_url FROM videos WHERE id = $1 AND user_id = $2`,
		videoID, userID,
	).Scan(&videoURL)
	if err != nil {
		http.Error(w, "動画が見つかりません", http.StatusNotFound)
		return
	}
	if videoURL == "" {
		http.Error(w, "動画ファイルがアップロードされていません", http.StatusBadRequest)
		return
	}

	// FFmpegに渡すURLは presigned GET（非公開バケットでも署名付きで読める）
	inputURL, err := h.minio.PresignedGetURL(r.Context(), videoURL)
	if err != nil {
		log.Printf("presigned get error: %v", err)
		http.Error(w, "サムネイル生成に失敗しました", http.StatusInternalServerError)
		return
	}
	objectName := fmt.Sprintf("thumbnails/%d/%d_auto.jpg", videoID, time.Now().Unix())

	thumbnailURL, err := h.minio.GenerateThumbnail(r.Context(), inputURL, objectName)
	if err != nil {
		log.Printf("thumbnail generation error: %v", err)
		http.Error(w, "サムネイル生成に失敗しました", http.StatusInternalServerError)
		return
	}

	// DBのthumbnail_urlを更新
	_, err = h.db.ExecContext(r.Context(),
		`UPDATE videos SET thumbnail_url = $1, updated_at = NOW() WHERE id = $2`,
		thumbnailURL, videoID,
	)
	if err != nil {
		http.Error(w, "DB更新に失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"thumbnail_url": h.resolver.PublicURL(thumbnailURL)})
}
