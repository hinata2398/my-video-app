# 既存の default VPC を参照（作らない・読むだけ）
data "aws_vpc" "default" {
  default = true
}

# その VPC 内のサブネット一覧を参照（ALB/ASG はマルチAZに置く）
data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}