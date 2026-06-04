package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hinata2398/my-video-app/backend/infrastructure/middleware"
	"github.com/hinata2398/my-video-app/backend/infrastructure/queue"
)

type TranscodeHandler struct {
	queue *queue.TranscodeQueue
	db    *sql.DB
}

func NewTranscodeHandler(q *queue.TranscodeQueue, db *sql.DB) *TranscodeHandler {
	return &TranscodeHandler{queue: q, db: db}
}

// Enqueue は動画をキューに追加して即座に202を返す
func (h *TranscodeHandler) Enqueue(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)
	videoID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

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

	// MinIO内部URLに変換
	internalURL := fmt.Sprintf("http://minio:9000/%s", extractPath(videoURL))
	objectName := fmt.Sprintf("videos/%d/%d_transcoded.mp4", videoID, time.Now().Unix())

	h.queue.Enqueue(queue.TranscodeJob{
		VideoID:    videoID,
		VideoURL:   internalURL,
		ObjectName: objectName,
	})

	// ステータスをpendingに更新
	h.db.ExecContext(r.Context(),
		`UPDATE videos SET status = 'pending', updated_at = NOW() WHERE id = $1`,
		videoID,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "pending"})
}

// Status はトランスコードのステータスを返す（フロントエンドがポーリングする）
func (h *TranscodeHandler) Status(w http.ResponseWriter, r *http.Request) {
	videoID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var status, message string
	err = h.db.QueryRowContext(r.Context(),
		`SELECT status, status_message FROM videos WHERE id = $1`,
		videoID,
	).Scan(&status, &message)
	if err != nil {
		http.Error(w, "動画が見つかりません", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": status, "message": message})
}

// extractPath は "http://localhost:9000/videos/foo/bar.mp4" → "videos/foo/bar.mp4"
func extractPath(publicURL string) string {
	const prefix = "http://localhost:9000/"
	if len(publicURL) > len(prefix) {
		return publicURL[len(prefix):]
	}
	return publicURL
}
