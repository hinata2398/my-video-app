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
	// 内部操作用（バケット管理など）
	internalClient *minio.Client
	// 署名付きURL生成用（ブラウザからアクセスできるホスト）
	publicClient   *minio.Client
	bucket         string
	publicEndpoint string
}

func NewMinioClient() (*MinioClient, error) {
	internalEndpoint := os.Getenv("MINIO_ENDPOINT")      // minio:9000
	publicEndpoint := os.Getenv("MINIO_PUBLIC_ENDPOINT") // localhost:9000
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	bucket := os.Getenv("MINIO_BUCKET")

	// Dockerネットワーク内部用クライアント
	internalClient, err := minio.New(internalEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	// 署名付きURL生成用クライアント（localhost:9000 で署名）
	publicClient, err := minio.New(publicEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	// バケットがなければ作成
	ctx := context.Background()
	exists, err := internalClient.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := internalClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
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
		if err := internalClient.SetBucketPolicy(ctx, bucket, policy); err != nil {
			return nil, err
		}
	}

	return &MinioClient{
		internalClient: internalClient,
		publicClient:   publicClient,
		bucket:         bucket,
		publicEndpoint: publicEndpoint,
	}, nil
}

func (m *MinioClient) PresignedUploadURL(ctx context.Context, objectName string) (string, error) {
	// publicClient で署名 → URLのホストが localhost:9000 になる
	u, err := m.publicClient.PresignedPutObject(ctx, m.bucket, objectName, 15*time.Minute)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func (m *MinioClient) PublicURL(objectName string) string {
	return fmt.Sprintf("http://%s/%s/%s", m.publicEndpoint, m.bucket, objectName)
}
