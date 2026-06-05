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

type Resolution struct {
	Name      string // "1080p", "720p", "480p"
	Height    int
	VideoBitrate string
	AudioBitrate string
	Bandwidth int // bits/sec（マスタープレイリスト用）
}

var Resolutions = []Resolution{
	{"1080p", 1080, "4000k", "192k", 4200000},
	{"720p",  720,  "2500k", "128k", 2700000},
	{"480p",  480,  "1000k", "96k",  1100000},
}

// GenerateHLS は動画を複数解像度のHLSに変換してMinIOに保存し、master.m3u8のURLを返す
func (m *MinioClient) GenerateHLS(ctx context.Context, videoURL, hlsPrefix string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "hls-*")
	if err != nil {
		return "", fmt.Errorf("tmp dir error: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// 元動画の高さを取得して、実際に生成する解像度を絞る
	actualHeight, err := getVideoHeight(ctx, videoURL)
	if err != nil {
		actualHeight = 9999 // 取得失敗時は全解像度生成
	}

	var generated []Resolution
	for _, r := range Resolutions {
		if r.Height > actualHeight {
			continue // 元動画より高解像度はスキップ
		}
		if err := m.generateResolution(ctx, videoURL, tmpDir, hlsPrefix, r); err != nil {
			return "", fmt.Errorf("transcode %s failed: %w", r.Name, err)
		}
		generated = append(generated, r)
	}

	// 1つも生成されなかった場合は最低解像度で生成
	if len(generated) == 0 {
		r := Resolutions[len(Resolutions)-1]
		if err := m.generateResolution(ctx, videoURL, tmpDir, hlsPrefix, r); err != nil {
			return "", err
		}
		generated = append(generated, r)
	}

	// マスタープレイリストを生成・アップロード
	masterURL, err := m.uploadMasterPlaylist(ctx, hlsPrefix, generated)
	if err != nil {
		return "", err
	}

	return masterURL, nil
}

func getVideoHeight(ctx context.Context, videoURL string) (int, error) {
	cmd := exec.CommandContext(ctx, "ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=height",
		"-of", "default=noprint_wrappers=1:nokey=1",
		videoURL,
	)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return 0, err
	}
	var h int
	fmt.Sscanf(strings.TrimSpace(out.String()), "%d", &h)
	return h, nil
}

func (m *MinioClient) generateResolution(ctx context.Context, videoURL, tmpDir, hlsPrefix string, r Resolution) error {
	resDir := filepath.Join(tmpDir, r.Name)
	if err := os.MkdirAll(resDir, 0755); err != nil {
		return err
	}

	playlistPath := filepath.Join(resDir, "index.m3u8")
	segmentPattern := filepath.Join(resDir, "segment_%03d.ts")

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", videoURL,
		"-c:v", "libx264",
		"-preset", "fast",
		"-b:v", r.VideoBitrate,
		// アスペクト比を維持しながら指定高さにスケール
		"-vf", fmt.Sprintf("scale=-2:%d", r.Height),
		"-c:a", "aac",
		"-b:a", r.AudioBitrate,
		"-hls_time", "6",
		"-hls_playlist_type", "vod",
		"-hls_segment_filename", segmentPattern,
		"-y",
		playlistPath,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg error: %w\nstderr: %s", err, stderr.String())
	}

	// 生成されたファイルをMinIOにアップロード
	return m.uploadDir(ctx, resDir, fmt.Sprintf("%s/%s", hlsPrefix, r.Name))
}

func (m *MinioClient) uploadDir(ctx context.Context, dir, objectPrefix string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read dir error: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		localPath := filepath.Join(dir, entry.Name())
		objectName := fmt.Sprintf("%s/%s", objectPrefix, entry.Name())

		f, err := os.Open(localPath)
		if err != nil {
			return err
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

func (m *MinioClient) uploadMasterPlaylist(ctx context.Context, hlsPrefix string, resolutions []Resolution) (string, error) {
	var sb strings.Builder
	sb.WriteString("#EXTM3U\n")
	sb.WriteString("#EXT-X-VERSION:3\n\n")

	for _, r := range resolutions {
		sb.WriteString(fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=x%d\n",
			r.Bandwidth, r.Height))
		// プレイヤーから見た相対パス
		sb.WriteString(fmt.Sprintf("%s/index.m3u8\n", r.Name))
	}

	content := sb.String()
	objectName := fmt.Sprintf("%s/master.m3u8", hlsPrefix)

	_, err := m.internalClient.PutObject(ctx, m.bucket, objectName,
		strings.NewReader(content), int64(len(content)),
		minio.PutObjectOptions{ContentType: "application/x-mpegURL"},
	)
	if err != nil {
		return "", fmt.Errorf("upload master playlist error: %w", err)
	}

	return m.PublicURL(objectName), nil
}
