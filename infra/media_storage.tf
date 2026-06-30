# メディアバケットの CORS（バケット本体は TF 管理外、CORSだけ管理）
resource "aws_s3_bucket_cors_configuration" "media" {
  bucket = "my-video-app-ryo-2026"

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT"]
    allowed_origins = [
      "https://${aws_cloudfront_distribution.frontend.domain_name}",
      "http://localhost:3000",
    ]
    expose_headers = ["ETag"]
  }
}

# backend IAM ユーザーの inline policy（オブジェクト read/write）
resource "aws_iam_user_policy" "backend_s3" {
  name = "s3-object-readwrite"
  user = "my-video-app-backend"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Sid      = "S3ObjectReadWrite"
      Effect   = "Allow"
      Action   = ["s3:GetObject", "s3:PutObject"]
      Resource = "arn:aws:s3:::my-video-app-ryo-2026/*"
    }]
  })
}
