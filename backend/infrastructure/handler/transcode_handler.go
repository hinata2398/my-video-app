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

type TranscodeHandler struct {
	minio *storage.MinioClient
	db    *sql.DB
}

func NewTranscodeHandler(minio *storage.MinioClient, db *sql.DB) *TranscodeHandler {
	return &TranscodeHandler{minio: minio, db: db}
}

func (h *TranscodeHandler) Transcode(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)
	videoID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// 元動画URLをDBから取得（自分の動画のみ）
	var originalURL string
	err = h.db.QueryRowContext(r.Context(),
		`SELECT video_url FROM videos WHERE id = $1 AND user_id = $2`,
		videoID, userID,
	).Scan(&originalURL)
	if err != nil {
		http.Error(w, "動画が見つかりません", http.StatusNotFound)
		return
	}
	if originalURL == "" {
		http.Error(w, "動画ファイルがアップロードされていません", http.StatusBadRequest)
		return
	}

	// FFmpegにはMinIO内部アドレスを渡す
	internalURL := fmt.Sprintf("http://minio:9000/%s", extractPath(originalURL))
	objectName := fmt.Sprintf("videos/%d/%d_transcoded.mp4", videoID, time.Now().Unix())

	log.Printf("transcoding video %d: %s", videoID, internalURL)

	transcodedURL, err := h.minio.Transcode(r.Context(), internalURL, objectName)
	if err != nil {
		log.Printf("transcode error video %d: %v", videoID, err)
		http.Error(w, "トランスコードに失敗しました", http.StatusInternalServerError)
		return
	}

	// video_url を変換済みファイルに更新
	_, err = h.db.ExecContext(r.Context(),
		`UPDATE videos SET video_url = $1, updated_at = NOW() WHERE id = $2`,
		transcodedURL, videoID,
	)
	if err != nil {
		http.Error(w, "DB更新に失敗しました", http.StatusInternalServerError)
		return
	}

	log.Printf("transcoding done video %d: %s", videoID, transcodedURL)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"video_url": transcodedURL})
}

// "http://localhost:9000/videos/foo/bar.mp4" → "videos/foo/bar.mp4"
func extractPath(publicURL string) string {
	const prefix = "http://localhost:9000/"
	if len(publicURL) > len(prefix) {
		return publicURL[len(prefix):]
	}
	return publicURL
}
