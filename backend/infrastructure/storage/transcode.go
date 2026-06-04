package storage

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/minio/minio-go/v7"
)

// Transcode は動画をH.264/AACのMP4に変換してMinIOに保存し、公開URLを返す
func (m *MinioClient) Transcode(ctx context.Context, videoURL, objectName string) (string, error) {
	// 一時ファイルに出力（faststartはseekが必要なためstdoutへのpipeが使えない）
	tmpFile, err := os.CreateTemp("", "transcode-*.mp4")
	if err != nil {
		return "", fmt.Errorf("tmp file error: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", videoURL,
		"-c:v", "libx264",        // 映像: H.264
		"-preset", "fast",        // エンコード速度優先
		"-crf", "23",             // 品質（18=高品質, 28=低品質）
		"-c:a", "aac",            // 音声: AAC
		"-b:a", "128k",
		"-movflags", "+faststart", // moovアトムを先頭に（ブラウザでの即再生に必須）
		"-y",                      // 上書き許可
		tmpFile.Name(),
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg error: %w\nstderr: %s", err, stderr.String())
	}

	// 変換済みファイルをMinIOにアップロード
	f, err := os.Open(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("open tmp file error: %w", err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return "", err
	}

	_, err = m.internalClient.PutObject(ctx, m.bucket, objectName, f, stat.Size(), minio.PutObjectOptions{
		ContentType: "video/mp4",
	})
	if err != nil {
		return "", fmt.Errorf("minio upload error: %w", err)
	}

	return m.PublicURL(objectName), nil
}
