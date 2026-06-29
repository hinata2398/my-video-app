terraform {
  required_version = ">= 1.5"

  backend "s3" {
    bucket       = "hinata-movie-tfstate"
    key          = "infra/terraform.tfstate" # バケット内での state の置き場所/名前
    region       = "ap-northeast-1"
    profile      = "mva-tf"
    encrypt      = true # state オブジェクトを暗号化
    use_lockfile = true # S3ネイティブロック（DynamoDB不要）
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}
