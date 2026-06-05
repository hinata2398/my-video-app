package storage

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
)

// GenerateHLS は動画をHLS形式（.m3u8 + .tsセグメント）に変換してMinIOに保存し、m3u8のURLを返す
func (m *MinioClient) GenerateHLS(ctx context.Context, videoURL, hlsPrefix string) (string, error) {
	// 一時ディレクトリに出力
	tmpDir, err := os.MkdirTemp("", "hls-*")
	if err != nil {
		return "", fmt.Errorf("tmp dir error: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	playlistPath := filepath.Join(tmpDir, "index.m3u8")
	segmentPattern := filepath.Join(tmpDir, "segment_%03d.ts")

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", videoURL,
		"-c:v", "libx264",
		"-preset", "fast",
		"-crf", "23",
		"-c:a", "aac",
		"-b:a", "128k",
		"-hls_time", "6",              // セグメント長（秒）
		"-hls_playlist_type", "vod",   // VOD（録画済み）形式
		"-hls_segment_filename", segmentPattern,
		"-y",
		playlistPath,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg error: %w\nstderr: %s", err, stderr.String())
	}

	// 生成されたファイルをMinIOにアップロード
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		return "", fmt.Errorf("read dir error: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		localPath := filepath.Join(tmpDir, entry.Name())
		objectName := fmt.Sprintf("%s/%s", hlsPrefix, entry.Name())

		f, err := os.Open(localPath)
		if err != nil {
			return "", fmt.Errorf("open file error: %w", err)
		}

		stat, _ := f.Stat()
		contentType := "video/MP2T"
		if strings.HasSuffix(entry.Name(), ".m3u8") {
			contentType = "application/x-mpegURL"
		}

		_, err = m.internalClient.PutObject(ctx, m.bucket, objectName, f, stat.Size(), minio.PutObjectOptions{
			ContentType: contentType,
		})
		f.Close()
		if err != nil {
			return "", fmt.Errorf("minio upload error (%s): %w", entry.Name(), err)
		}
	}

	return m.PublicURL(fmt.Sprintf("%s/index.m3u8", hlsPrefix)), nil
}
