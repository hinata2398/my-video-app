variable "region" {
  type    = string
  default = "ap-northeast-1"
}

variable "aws_profile" {
  type    = string
  default = "mva-tf"
}

variable "alert_email" {}

variable "db_password" {
  type      = string
  sensitive = true # plan/apply のログに出さない
}
