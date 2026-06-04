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
	VideoID    int64
	VideoURL   string // MinIO内部URL
	ObjectName string // 保存先オブジェクト名
}

type TranscodeQueue struct {
	jobs  chan TranscodeJob
	minio *storage.MinioClient
	db    *sql.DB
}

func NewTranscodeQueue(minio *storage.MinioClient, db *sql.DB, workers int) *TranscodeQueue {
	q := &TranscodeQueue{
		jobs:  make(chan TranscodeJob, 100), // バッファ100件
		minio: minio,
		db:    db,
	}
	// workerをgoroutineで起動
	for i := 0; i < workers; i++ {
		go q.worker(i)
	}
	return q
}

// Enqueue はジョブをキューに追加して即座に返る
func (q *TranscodeQueue) Enqueue(job TranscodeJob) {
	q.jobs <- job
}

// worker はキューからジョブを取り出して処理するgoroutine
func (q *TranscodeQueue) worker(id int) {
	log.Printf("transcode worker %d started", id)
	for job := range q.jobs {
		q.process(job)
	}
}

func (q *TranscodeQueue) process(job TranscodeJob) {
	log.Printf("transcode start: video_id=%d", job.VideoID)

	// ステータスを processing に更新
	q.updateStatus(job.VideoID, "processing", "")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// トランスコード実行
	transcodedURL, err := q.minio.Transcode(ctx, job.VideoURL, job.ObjectName)
	if err != nil {
		log.Printf("transcode failed: video_id=%d, err=%v", job.VideoID, err)
		q.updateStatus(job.VideoID, "error", err.Error())
		return
	}

	// video_url を変換済みに更新
	_, err = q.db.Exec(
		`UPDATE videos SET video_url = $1, updated_at = NOW() WHERE id = $2`,
		transcodedURL, job.VideoID,
	)
	if err != nil {
		log.Printf("db update failed: video_id=%d, err=%v", job.VideoID, err)
		q.updateStatus(job.VideoID, "error", err.Error())
		return
	}

	// サムネイル自動生成
	thumbObjectName := fmt.Sprintf("thumbnails/%d/%d_auto.jpg", job.VideoID, time.Now().Unix())
	thumbnailURL, err := q.minio.GenerateThumbnail(ctx, job.VideoURL, thumbObjectName)
	if err != nil {
		log.Printf("thumbnail failed (non-fatal): video_id=%d, err=%v", job.VideoID, err)
		// サムネイル失敗はdoneにする（動画再生には影響しない）
	} else {
		q.db.Exec(
			`UPDATE videos SET thumbnail_url = $1, updated_at = NOW() WHERE id = $2 AND thumbnail_url = ''`,
			thumbnailURL, job.VideoID,
		)
	}

	q.updateStatus(job.VideoID, "done", "")
	log.Printf("transcode done: video_id=%d url=%s", job.VideoID, transcodedURL)
}

func (q *TranscodeQueue) updateStatus(videoID int64, status, errMsg string) {
	q.db.Exec(
		`UPDATE videos SET status = $1, status_message = $2, updated_at = NOW() WHERE id = $3`,
		status, errMsg, videoID,
	)
}
