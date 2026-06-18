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
	internalClient   *minio.Client
	publicClient     *minio.Client
	bucket           string
	publicEndpoint   string
	internalEndpoint string
}

func NewMinioClient() (*MinioClient, error) {
	internalEndpoint := os.Getenv("MINIO_ENDPOINT")      // minio:9000
	publicEndpoint := os.Getenv("MINIO_PUBLIC_ENDPOINT") // localhost:9000
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	bucket := os.Getenv("MINIO_BUCKET")
	secure := os.Getenv("MINIO_SECURE") == "true" // 既定 false
	region := os.Getenv("MINIO_REGION")
	if region == "" {
		region = "us-east-1" // 既定
	}
	autoCreate := os.Getenv("MINIO_AUTO_CREATE_BUCKET") != "false" // 既定 true

	// 内部操作用（バケット管理）
	internalClient, err := minio.New(internalEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
		Region: region,
	})
	if err != nil {
		return nil, err
	}

	// 署名付きURL生成用（localhost:9000 で署名、Region固定でHTTP通信不要）
	publicClient, err := minio.New(publicEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
		Region: region, // Region固定でGetBucketLocationのHTTPリクエストを省略
	})
	if err != nil {
		return nil, err
	}

	if autoCreate {
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
	}

	return &MinioClient{
		internalClient:   internalClient,
		publicClient:     publicClient,
		bucket:           bucket,
		publicEndpoint:   publicEndpoint,
		internalEndpoint: internalEndpoint,
	}, nil
}

func (m *MinioClient) PresignedUploadURL(ctx context.Context, objectName string) (string, error) {
	// publicClient で署名 → Host が localhost:9000 で署名される
	// Region固定により GetBucketLocation の HTTP通信が発生しない
	u, err := m.publicClient.PresignedPutObject(ctx, m.bucket, objectName, 15*time.Minute)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func (m *MinioClient) PresignedGetURL(ctx context.Context, key string) (string, error) {
	u, err := m.internalClient.PresignedGetObject(ctx, m.bucket, key, 1*time.Hour, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func (m *MinioClient) GetObject(ctx context.Context, key string) (*minio.Object, error) {
	return m.internalClient.GetObject(ctx, m.bucket, key, minio.GetObjectOptions{})
}
