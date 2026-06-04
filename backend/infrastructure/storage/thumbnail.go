package storage

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/minio/minio-go/v7"
)

// GenerateThumbnail は動画URLからFFmpegでサムネイルを生成し、MinIOに保存してURLを返す
func (m *MinioClient) GenerateThumbnail(ctx context.Context, videoURL, objectName string) (string, error) {
	// FFmpegで動画の1秒目のフレームをJPEGとして標準出力に出力
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", videoURL,
		"-ss", "00:00:01",
		"-vframes", "1",
		"-vf", "scale=480:-1",
		"-f", "image2pipe",
		"-vcodec", "mjpeg",
		"-q:v", "3",
		"pipe:1",
	)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg error: %w, stderr: %s", err, stderr.String())
	}

	// MinIOにアップロード
	data := bytes.NewReader(out.Bytes())
	_, err := m.internalClient.PutObject(ctx, m.bucket, objectName, data, int64(out.Len()), minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		return "", fmt.Errorf("minio upload error: %w", err)
	}

	return m.PublicURL(objectName), nil
}
