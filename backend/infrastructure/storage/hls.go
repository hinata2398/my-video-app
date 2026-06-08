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

// 変換する画質の種類を定義
type qualityProfile struct {
	name       string // フォルダ名（例: "360p"）
	resolution string // 解像度（例: "640:360"）
	bitrate    string // 映像の品質（例: "800k"）
}

var qualityProfiles = []qualityProfile{
	{name: "360p", resolution: "640:360", bitrate: "800k"},
	{name: "720p", resolution: "1280:720", bitrate: "2800k"},
}

// GenerateHLS は動画をHLS形式（.m3u8 + .tsセグメント）に変換してMinIOに保存し、m3u8のURLを返す
func (m *MinioClient) GenerateHLS(ctx context.Context, videoURL, hlsPrefix string) (string, error) {
	// 一時ディレクトリに出力
	tmpDir, err := os.MkdirTemp("", "hls-*")
	if err != nil {
		return "", fmt.Errorf("tmp dir error: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// 画質ごとにFFmpegを実行
	for _, q := range qualityProfiles {
		qualityDir := filepath.Join(tmpDir, q.name)
		if err := os.MkdirAll(qualityDir, 0755); err != nil {
			return "", fmt.Errorf("mkdir error (%s): %w", q.name, err)
		}

		playlistPath := filepath.Join(qualityDir, "index.m3u8")
		segmentPattern := filepath.Join(qualityDir, "segment_%03d.ts")

		cmd := exec.CommandContext(ctx, "ffmpeg",
			"-i", videoURL,
			"-c:v", "libx264",
			"-preset", "fast",
			"-vf", "scale="+q.resolution,
			"-b:v", q.bitrate,
			"-c:a", "aac",
			"-b:a", "128k",
			"-hls_time", "6",
			"-hls_playlist_type", "vod",
			"-hls_segment_filename", segmentPattern,
			"-y",
			playlistPath,
		)

		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("ffmpeg error (%s): %w\nstderr: %s", q.name, err, stderr.String())
		}
	}

	// master.m3u8 を作成
	masterPath := filepath.Join(tmpDir, "master.m3u8")
	masterContent := "#EXTM3U\n"
	for _, q := range qualityProfiles {
		masterContent += fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%s,RESOLUTION=%s\n%s/index.m3u8\n",
			bandwidthFromBitrate(q.bitrate),
			resolutionLabel(q.resolution),
			q.name,
		)
	}
	if err := os.WriteFile(masterPath, []byte(masterContent), 0644); err != nil {
		return "", fmt.Errorf("master.m3u8 write error: %w", err)
	}

	// 生成されたファイルをMinIOにアップロード（サブフォルダも含む）
	err = uploadDir(ctx, m, tmpDir, hlsPrefix)
	if err != nil {
		return "", err
	}

	return m.PublicURL(fmt.Sprintf("%s/master.m3u8", hlsPrefix)), nil
}

// "800k" → "800000" に変換（master.m3u8 に必要な単位）
func bandwidthFromBitrate(bitrate string) string {
	bitrate = strings.TrimSuffix(bitrate, "k")
	val := 0
	fmt.Sscanf(bitrate, "%d", &val)
	return fmt.Sprintf("%d", val*1000)
}

// "640:360" → "640x360" に変換（master.m3u8 の表記ルール）
func resolutionLabel(resolution string) string {
	return strings.ReplaceAll(resolution, ":", "x")
}

// tmpDir 以下のファイルを再帰的にMinIOへアップロードする
func uploadDir(ctx context.Context, m *MinioClient, localDir, prefix string) error {
	entries, err := os.ReadDir(localDir)
	if err != nil {
		return fmt.Errorf("read dir error: %w", err)
	}

	for _, entry := range entries {
		localPath := filepath.Join(localDir, entry.Name())
		objectName := fmt.Sprintf("%s/%s", prefix, entry.Name())

		if entry.IsDir() {
			// サブフォルダは再帰的に処理
			if err := uploadDir(ctx, m, localPath, objectName); err != nil {
				return err
			}
			continue
		}

		f, err := os.Open(localPath)
		if err != nil {
			return fmt.Errorf("open file error: %w", err)
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
			return fmt.Errorf("minio upload error (%s): %w", entry.Name(), err)
		}
	}
	return nil
}
