resource "aws_route53_zone" "main" {
  name = "hinata-movie.net"
}

# お名前.comに設定するNS 4つを出力
output "route53_nameservers" {
  value = aws_route53_zone.main.name_servers
}

# 証明書をリクエスト(ALB用はALBと同じ ap-northeast-1 で発行)
resource "aws_acm_certificate" "main" {
  domain_name       = "hinata-movie.net"
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

# ACMが要求する検証用CNAMEをRoute53に自動作成
resource "aws_route53_record" "cert_validation" {
  for_each = {
    for dvo in aws_acm_certificate.main.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      type   = dvo.resource_record_type
      record = dvo.resource_record_value
    }
  }

  zone_id         = aws_route53_zone.main.zone_id
  name            = each.value.name
  type            = each.value.type
  records         = [each.value.record]
  ttl             = 60
  allow_overwrite = true
}

# 検証完了まで待つ(これ自体はAWSリソースではなく"待機"用)
resource "aws_acm_certificate_validation" "main" {
  certificate_arn         = aws_acm_certificate.main.arn
  validation_record_fqdns = [for r in aws_route53_record.cert_validation : r.fqdn]
}

resource "aws_route53_record" "apex" {
  zone_id = aws_route53_zone.main.zone_id
  name    = "hinata-movie.net"
  type    = "A"

  alias {
    name                   = aws_lb.app.dns_name
    zone_id                = aws_lb.app.zone_id
    evaluate_target_health = true
  }
}