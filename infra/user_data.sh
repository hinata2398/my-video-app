#!/bin/bash
#dockerとgitをインストール
dnf install -y docker git amazon-efs-utils
#dockerをシステムの起動時に自動で開始・実行
systemctl enable --now docker
#gitリポジトリのクローン
git clone https://github.com/hinata2398/my-video-app.git
cd /my-video-app
#dockerイメージの作成
docker build -t mva-backend ./backend

mkdir -p /mnt/efs/transcode-scratch
mount -t efs -o tls fs-04c876346cf22461d:/ /mnt/efs/transcode-scratch

#.envファイルの作成
cat > app.env <<EOF
DATABASE_URL=$(aws ssm get-parameter --name /my-video-app/DATABASE_URL --with-decryption --query Parameter.Value --output text --region ap-northeast-1)
JWT_SECRET=$(aws ssm get-parameter --name /my-video-app/JWT_SECRET --with-decryption --query Parameter.Value --output text --region ap-northeast-1)
MINIO_ACCESS_KEY=$(aws ssm get-parameter --name /my-video-app/MINIO_ACCESS_KEY --with-decryption --query Parameter.Value --output text --region ap-northeast-1)
MINIO_SECRET_KEY=$(aws ssm get-parameter --name /my-video-app/MINIO_SECRET_KEY --with-decryption --query Parameter.Value --output text --region ap-northeast-1)
MINIO_ENDPOINT=s3.ap-northeast-1.amazonaws.com
MINIO_PUBLIC_ENDPOINT=s3.ap-northeast-1.amazonaws.com
MINIO_BUCKET=my-video-app-ryo-2026
MINIO_SECURE=true
MINIO_REGION=ap-northeast-1
MINIO_AUTO_CREATE_BUCKET=false
MEDIA_DELIVERY=cloudfront
MEDIA_BASE_URL=https://d2acjx9xgv33qc.cloudfront.net
EOF

docker run -d --name backend -p 8080:8080 --env-file app.env \
  -v /mnt/efs/transcode-scratch:/tmp/transcode \
  --log-driver=awslogs \
  --log-opt awslogs-region=ap-northeast-1 \
  --log-opt awslogs-group=/my-video-app/backend \
  --log-opt awslogs-create-group=false \
  mva-backend