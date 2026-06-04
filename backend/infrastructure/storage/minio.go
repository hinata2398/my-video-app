package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	client         *minio.Client
	bucket         string
	publicEndpoint string
}

func NewMinioClient() (*MinioClient, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	bucket := os.Getenv("MINIO_BUCKET")
	publicEndpoint := os.Getenv("MINIO_PUBLIC_ENDPOINT")

	client, err := minio.New(endpoint, &minio.Options{
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
		// パブリック読み取りポリシーを設定
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
		client:         client,
		bucket:         bucket,
		publicEndpoint: publicEndpoint,
	}, nil
}

func (m *MinioClient) PresignedUploadURL(ctx context.Context, objectName string) (string, error) {
	u, err := m.client.PresignedPutObject(ctx, m.bucket, objectName, 15*time.Minute)
	if err != nil {
		return "", err
	}
	// Docker内部ホスト(minio:9000)をブラウザからアクセス可能なホスト(localhost:9000)に置換
	u.Host = m.publicEndpoint
	u.Scheme = "http"
	return u.String(), nil
}

func (m *MinioClient) PublicURL(objectName string) string {
	return fmt.Sprintf("http://%s/%s/%s", m.publicEndpoint, m.bucket, objectName)
}
