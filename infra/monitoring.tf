# 通知先トピック
resource "aws_sns_topic" "alerts" {
  name = "my-video-app-alerts"
}

# メール購読(★applyするとAWSから確認メールが来る→リンク踏むまで PendingConfirmation)
resource "aws_sns_topic_subscription" "alerts_email" {
  topic_arn = aws_sns_topic.alerts.arn
  protocol  = "email"
  endpoint  = var.alert_email
}

# 最初のアラーム:ALBの正常ターゲット数が2未満になったら
resource "aws_cloudwatch_metric_alarm" "alb_unhealthy" {
  alarm_name          = "my-video-app-alb-unhealthy-hosts"
  namespace           = "AWS/ApplicationELB"
  metric_name         = "HealthyHostCount"
  statistic           = "Minimum"
  comparison_operator = "LessThanThreshold"
  threshold           = 2
  period              = 60
  evaluation_periods  = 2

  dimensions = {
    TargetGroup  = aws_lb_target_group.app.arn_suffix # ★arn ではなく arn_suffix
    LoadBalancer = aws_lb.app.arn_suffix              # ★同上
  }

  alarm_actions = [aws_sns_topic.alerts.arn]
  ok_actions    = [aws_sns_topic.alerts.arn]
}

# RDS: 空きストレージが2GB未満(★単位はバイト)
resource "aws_cloudwatch_metric_alarm" "rds_low_storage" {
  alarm_name          = "my-video-app-rds-low-storage"
  namespace           = "AWS/RDS"
  metric_name         = "FreeStorageSpace"
  statistic           = "Average"
  comparison_operator = "LessThanThreshold"
  threshold           = 2000000000 # 2GB(バイト)
  period              = 300
  evaluation_periods  = 1
  dimensions          = { DBInstanceIdentifier = aws_db_instance.main.identifier }
  alarm_actions       = [aws_sns_topic.alerts.arn]
  ok_actions          = [aws_sns_topic.alerts.arn]
  treat_missing_data  = "notBreaching" # ★RDS停止中はデータ欠損→誤発報させない
}

# RDS: CPU使用率が80%超
resource "aws_cloudwatch_metric_alarm" "rds_high_cpu" {
  alarm_name          = "my-video-app-rds-high-cpu"
  namespace           = "AWS/RDS"
  metric_name         = "CPUUtilization"
  statistic           = "Average"
  comparison_operator = "GreaterThanThreshold"
  threshold           = 80
  period              = 300
  evaluation_periods  = 2
  dimensions          = { DBInstanceIdentifier = aws_db_instance.main.identifier }
  alarm_actions       = [aws_sns_topic.alerts.arn]
  ok_actions          = [aws_sns_topic.alerts.arn]
  treat_missing_data  = "notBreaching"
}

# EC2(アプリ): ASG全体のCPU平均が80%超 ★namespaceはAWS/EC2、次元はASG名
resource "aws_cloudwatch_metric_alarm" "ec2_high_cpu" {
  alarm_name          = "my-video-app-ec2-high-cpu"
  namespace           = "AWS/EC2"
  metric_name         = "CPUUtilization"
  statistic           = "Average"
  comparison_operator = "GreaterThanThreshold"
  threshold           = 80
  period              = 300
  evaluation_periods  = 2
  dimensions          = { AutoScalingGroupName = aws_autoscaling_group.app.name }
  alarm_actions       = [aws_sns_topic.alerts.arn]
  ok_actions          = [aws_sns_topic.alerts.arn]
  treat_missing_data  = "notBreaching"
}

resource "aws_cloudwatch_dashboard" "main" {
  dashboard_name = "my-video-app"

  dashboard_body = jsonencode({
    widgets = [
      {
        type = "metric", x = 0, y = 0, width = 12, height = 6,
        properties = {
          title  = "ALB Healthy Hosts",
          region = var.region,
          view   = "timeSeries",
          metrics = [
            ["AWS/ApplicationELB", "HealthyHostCount",
              "TargetGroup", aws_lb_target_group.app.arn_suffix,
            "LoadBalancer", aws_lb.app.arn_suffix]
          ]
        }
      },
      {
        type = "metric", x = 12, y = 0, width = 12, height = 6,
        properties = {
          title  = "RDS CPU / FreeStorage",
          region = var.region,
          view   = "timeSeries",
          metrics = [
            ["AWS/RDS", "CPUUtilization", "DBInstanceIdentifier", aws_db_instance.main.identifier],
            ["AWS/RDS", "FreeStorageSpace", "DBInstanceIdentifier", aws_db_instance.main.identifier]
          ]
        }
      },
      {
        type = "metric", x = 0, y = 6, width = 12, height = 6,
        properties = {
          title  = "EC2 (ASG) CPU",
          region = var.region,
          view   = "timeSeries",
          metrics = [
            ["AWS/EC2", "CPUUtilization", "AutoScalingGroupName", aws_autoscaling_group.app.name]
          ]
        }
      }
    ]
  })
}

resource "aws_cloudwatch_log_group" "app" {
  name              = "/my-video-app/backend"
  retention_in_days = 14 # ★無期限保持を避ける(課金&肥大対策)
}

resource "aws_iam_role_policy" "ec2_logs" {
  name = "cloudwatch-logs-write"
  role = "ec2-ssm-role" # 既存ロールに名前で紐付け(本体はTF外のまま)

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = ["logs:CreateLogStream", "logs:PutLogEvents"]
      Resource = "${aws_cloudwatch_log_group.app.arn}:*"
    }]
  })
}