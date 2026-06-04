package storage

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	client         *minio.Client
	bucket         string
	internalEndpoint string
	publicEndpoint   string
}

func NewMinioClient() (*MinioClient, error) {
	internalEndpoint := os.Getenv("MINIO_ENDPOINT")      // minio:9000
	publicEndpoint := os.Getenv("MINIO_PUBLIC_ENDPOINT") // localhost:9000
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	bucket := os.Getenv("MINIO_BUCKET")

	client, err := minio.New(internalEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
		policy := fmt.Sprintf(`{
			"Version":"2012-10-17",
			"Statement":[{
				"Effect":"Allow",
				"Principal":"*",
				"Action":["s3:GetObject","s3:PutObject"],
				"Resource":["arn:aws:s3:::%s/*"]
			}]
		}`, bucket)
		if err := client.SetBucketPolicy(ctx, bucket, policy); err != nil {
			return nil, err
		}
	}

	return &MinioClient{
		client:           client,
		bucket:           bucket,
		internalEndpoint: internalEndpoint,
		publicEndpoint:   publicEndpoint,
	}, nil
}

func (m *MinioClient) PresignedUploadURL(ctx context.Context, objectName string) (string, error) {
	// 内部クライアントで署名付きURL生成（ネットワーク通信なし）
	u, err := m.client.PresignedPutObject(ctx, m.bucket, objectName, 15*time.Minute)
	if err != nil {
		return "", err
	}
	// minio:9000 → localhost:9000 に文字列置換
	// MinIOは受信リクエストのホストではなく認証情報で署名検証するため有効
	rawURL := u.String()
	rawURL = strings.Replace(rawURL, m.internalEndpoint, m.publicEndpoint, 1)
	return rawURL, nil
}

func (m *MinioClient) PublicURL(objectName string) string {
	return fmt.Sprintf("http://%s/%s/%s", m.publicEndpoint, m.bucket, objectName)
}
