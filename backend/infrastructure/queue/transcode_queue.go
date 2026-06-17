package queue

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/hinata2398/my-video-app/backend/infrastructure/storage"
)

type TranscodeJob struct {
	VideoID  int64
	VideoKey string // MinIO内部URL（http://minio:9000/...）
}

type TranscodeQueue struct {
	jobs  chan TranscodeJob
	minio *storage.MinioClient
	db    *sql.DB
}

func NewTranscodeQueue(minio *storage.MinioClient, db *sql.DB, workers int) *TranscodeQueue {
	q := &TranscodeQueue{
		jobs:  make(chan TranscodeJob, 100),
		minio: minio,
		db:    db,
	}
	for i := 0; i < workers; i++ {
		go q.worker(i)
	}
	return q
}

func (q *TranscodeQueue) Enqueue(job TranscodeJob) {
	q.jobs <- job
}

func (q *TranscodeQueue) worker(id int) {
	log.Printf("transcode worker %d started", id)
	for job := range q.jobs {
		q.process(job)
	}
}

func (q *TranscodeQueue) process(job TranscodeJob) {
	log.Printf("transcode start: video_id=%d", job.VideoID)
	q.updateStatus(job.VideoID, "processing", "")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// キー → ffmpegが読める内部URL（minio:9000）に変換
	internalURL := q.minio.InternalURL(job.VideoKey)

	// 1. HLS変換（.m3u8 + .tsセグメント）
	hlsPrefix := fmt.Sprintf("hls/%d/%d", job.VideoID, time.Now().Unix())
	hlsKey, err := q.minio.GenerateHLS(ctx, internalURL, hlsPrefix)
	if err != nil {
		log.Printf("HLS failed: video_id=%d, err=%v", job.VideoID, err)
		q.updateStatus(job.VideoID, "error", err.Error())
		return
	}

	// video_url をHLSプレイリストのキーに更新
	_, err = q.db.Exec(
		`UPDATE videos SET video_url = $1, updated_at = NOW() WHERE id = $2`,
		hlsKey, job.VideoID,
	)
	if err != nil {
		log.Printf("db update failed: video_id=%d, err=%v", job.VideoID, err)
		q.updateStatus(job.VideoID, "error", err.Error())
		return
	}

	// 2. サムネイル自動生成（thumbnail_urlが空の場合のみ）
	thumbObjectName := fmt.Sprintf("thumbnails/%d/%d_auto.jpg", job.VideoID, time.Now().Unix())
	thumbnailKey, err := q.minio.GenerateThumbnail(ctx, internalURL, thumbObjectName)
	if err != nil {
		log.Printf("thumbnail failed (non-fatal): video_id=%d, err=%v", job.VideoID, err)
	} else {
		q.db.Exec(
			`UPDATE videos SET thumbnail_url = $1, updated_at = NOW() WHERE id = $2 AND thumbnail_url = ''`,
			thumbnailKey, job.VideoID,
		)
	}

	q.updateStatus(job.VideoID, "done", "")
	log.Printf("transcode done: video_id=%d hls=%s", job.VideoID, hlsKey)
}

func (q *TranscodeQueue) updateStatus(videoID int64, status, errMsg string) {
	q.db.Exec(
		`UPDATE videos SET status = $1, status_message = $2, updated_at = NOW() WHERE id = $3`,
		status, errMsg, videoID,
	)
}
