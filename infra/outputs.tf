output "vpc_id" {
  value = aws_vpc.main.id
}

output "private_subnet_ids" {
  value = [aws_subnet.private_a.id, aws_subnet.private_c.id]
}

output "alb_dns_name" {
  value = aws_lb.app.dns_name
}
