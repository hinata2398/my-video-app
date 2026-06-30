package queue

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/hinata2398/my-video-app/backend/infrastructure/storage"
)

type TranscodeJob struct {
	VideoID  int64
	VideoKey string // MinIO内部URL（http://minio:9000/...）
}

type TranscodeQueue struct {
	sqs      *sqs.Client
	queueURL string
	minio    *storage.MinioClient
	db       *sql.DB
}

func NewTranscodeQueue(minio *storage.MinioClient, db *sql.DB, queueURL string) (*TranscodeQueue, error) {
	cfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	q := &TranscodeQueue{
		sqs: sqs.NewFromConfig(cfg), queueURL: queueURL, minio: minio, db: db,
	}
	q.StartConsumer()
	return q, nil
}

func (q *TranscodeQueue) Enqueue(ctx context.Context, job TranscodeJob) error {
	body, err := json.Marshal(job)
	if err != nil {
		return err
	}
	_, err = q.sqs.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(q.queueURL),
		MessageBody: aws.String(string(body)),
	})
	return err
}

func (q *TranscodeQueue) StartConsumer() { go q.poll() }

func (q *TranscodeQueue) poll() {
	log.Println("transcode consumer started (SQS polling)")
	for {
		out, err := q.sqs.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(q.queueURL),
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     20, // ロングポーリング
		})
		if err != nil {
			log.Printf("sqs receive error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		for _, m := range out.Messages {
			var job TranscodeJob
			if err := json.Unmarshal([]byte(*m.Body), &job); err != nil {
				log.Printf("bad message, dropping: %v", err)
				q.deleteMessage(m.ReceiptHandle) // 壊れたメッセージは捨てる
				continue
			}
			if err := q.process(job); err != nil {
				// 消さない → 可視性タイムアウト後に再配信 → 3回でDLQ
				log.Printf("process failed (retry/DLQ): video_id=%d err=%v", job.VideoID, err)
				continue
			}
			q.deleteMessage(m.ReceiptHandle) // 成功したら削除
		}
	}
}

func (q *TranscodeQueue) deleteMessage(rh *string) {
	_, err := q.sqs.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
		QueueUrl: aws.String(q.queueURL), ReceiptHandle: rh,
	})
	if err != nil {
		log.Printf("sqs delete error: %v", err)
	}
}

func (q *TranscodeQueue) process(job TranscodeJob) error {
	log.Printf("transcode start: video_id=%d", job.VideoID)
	q.updateStatus(job.VideoID, "processing", "")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// キー → ffmpegが読める内部URL（minio:9000）に変換
	inputURL, err := q.minio.PresignedGetURL(ctx, job.VideoKey)
	if err != nil {
		q.updateStatus(job.VideoID, "error", err.Error())
		return err
	}

	// 1. HLS変換（.m3u8 + .tsセグメント）
	hlsPrefix := fmt.Sprintf("hls/%d/%d", job.VideoID, time.Now().Unix())
	hlsKey, err := q.minio.GenerateHLS(ctx, inputURL, hlsPrefix)
	if err != nil {
		log.Printf("HLS failed: video_id=%d, err=%v", job.VideoID, err)
		q.updateStatus(job.VideoID, "error", err.Error())
		return err
	}

	// video_url をHLSプレイリストのキーに更新
	_, err = q.db.Exec(
		`UPDATE videos SET video_url = $1, updated_at = NOW() WHERE id = $2`,
		hlsKey, job.VideoID,
	)
	if err != nil {
		log.Printf("db update failed: video_id=%d, err=%v", job.VideoID, err)
		q.updateStatus(job.VideoID, "error", err.Error())
		return err
	}

	// 2. サムネイル自動生成（thumbnail_urlが空の場合のみ）
	thumbObjectName := fmt.Sprintf("thumbnails/%d/%d_auto.jpg", job.VideoID, time.Now().Unix())
	thumbnailKey, err := q.minio.GenerateThumbnail(ctx, inputURL, thumbObjectName)
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
	return nil
}

func (q *TranscodeQueue) updateStatus(videoID int64, status, errMsg string) {
	q.db.Exec(
		`UPDATE videos SET status = $1, status_message = $2, updated_at = NOW() WHERE id = $3`,
		status, errMsg, videoID,
	)
}
