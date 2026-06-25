output "default_vpc_id" {
  value = data.aws_vpc.default.id
}

output "default_subnet_ids" {
  value = data.aws_subnets.default.ids
}

output "alb_dns_name" {
  value = aws_lb.app.dns_name
}