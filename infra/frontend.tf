resource "aws_s3_bucket" "frontend" {
  bucket = "hinata-movie-frontend"
}

resource "aws_s3_bucket_public_access_block" "frontend" {
  bucket = aws_s3_bucket.frontend.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_cloudfront_origin_access_control" "frontend" {
  name                              = "hinata-movie-frontend-oac"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

resource "aws_cloudfront_distribution" "frontend" {
  enabled             = true         # 配信を有効化
  default_root_object = "index.html" # / にアクセスしたとき返すファイル

  origin {
    # S3のリージョナルドメイン（バケット名.s3.ap-northeast-1.amazonaws.com）
    domain_name = aws_s3_bucket.frontend.bucket_regional_domain_name

    # このoriginの識別名（後でcache_behaviorから参照する）
    origin_id = "s3-frontend"

    # さっき作ったOACを紐付ける
    origin_access_control_id = aws_cloudfront_origin_access_control.frontend.id
  }

  default_cache_behavior {
    target_origin_id       = "s3-frontend"       # どのoriginに転送するか
    viewer_protocol_policy = "redirect-to-https" # httpは301でhttpsへ
    allowed_methods        = ["GET", "HEAD"]     # 静的ファイルはGET/HEADのみ
    cached_methods         = ["GET", "HEAD"]     # キャッシュするメソッド

    forwarded_values {
      query_string = false         # クエリ文字列はS3に転送しない（キャッシュ効率UP）
      cookies { forward = "none" } # Cookieも転送しない
    }
  }

  restrictions {
    geo_restriction {
      restriction_type = "none" # 地域制限なし（全世界から配信）
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
    # *.cloudfront.net のデフォルト証明書を使う
    # カスタムドメイン(hinata-movie.net)にする場合はacm_certificate_arn を指定する
  }

  custom_error_response {
    error_code         = 404
    response_code      = 200
    response_page_path = "/index.html"
  }
}

# バケットポリシーをS3バケットに紐付ける
resource "aws_s3_bucket_policy" "frontend" {
  bucket = aws_s3_bucket.frontend.id

  # jsonencode = HCLのオブジェクトをJSON文字列に変換（monitoring.tfのdashboardと同じ手法）
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"

        # 誰に許可するか → CloudFrontサービス
        Principal = {
          Service = "cloudfront.amazonaws.com"
        }

        # 何を許可するか → オブジェクトの読み取りのみ
        Action = "s3:GetObject"

        # どのリソースに → バケット内の全オブジェクト
        Resource = "${aws_s3_bucket.frontend.arn}/*"

        # 条件: このCloudFrontディストリビューションからのリクエストのみ
        # → 他人がCloudFront経由でアクセスしてもブロック
        Condition = {
          StringEquals = {
            "AWS:SourceArn" = aws_cloudfront_distribution.frontend.arn
          }
        }
      }
    ]
  })
}
