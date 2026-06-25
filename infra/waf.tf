resource "aws_wafv2_web_acl" "main" {
  name  = "my-video-app-waf"
  scope = "REGIONAL"
  default_action {
    allow {}
  }
  visibility_config {
    cloudwatch_metrics_enabled = true
    metric_name                = "myVideoAppWaf" # 英数字のみ（ハイフン不可）。例: myVideoAppWaf
    sampled_requests_enabled   = true
  }
  rule {
    name     = "AWS-Common"
    priority = 1

    override_action {
      none {}
    }

    statement {
      managed_rule_group_statement {
        name        = "AWSManagedRulesCommonRuleSet"
        vendor_name = "AWS"
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name                = "awsCommon"
      sampled_requests_enabled   = true
    }
  }
}

resource "aws_wafv2_web_acl_association" "alb" {
  resource_arn = aws_lb.app.arn
  web_acl_arn  = aws_wafv2_web_acl.main.arn
}
