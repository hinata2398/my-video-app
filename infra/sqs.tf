# DLQ: 失敗し続けたメッセージの隔離先
resource "aws_sqs_queue" "transcode_dlq" {
  name                      = "my-video-app-transcode-dlq"
  message_retention_seconds = 1209600 # 14日（最大）= 調査猶予
}

# メイン: トランスコードジョブ
resource "aws_sqs_queue" "transcode" {
  name                       = "my-video-app-transcode"
  visibility_timeout_seconds = 1800  # ★ffmpeg最大処理時間以上に（30分）
  receive_wait_time_seconds  = 20    # ロングポーリング
  message_retention_seconds  = 86400 # 1日

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.transcode_dlq.arn
    maxReceiveCount     = 3 # 3回失敗で DLQ
  })
}

# 既存EC2ロール(ec2-ssm-role, TF管理外)に SQS 権限を名前で付与（5dのlogsと同じ手法）
resource "aws_iam_role_policy" "ec2_sqs" {
  name = "transcode-sqs"
  role = "ec2-ssm-role"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "sqs:SendMessage",
        "sqs:ReceiveMessage",
        "sqs:DeleteMessage",
        "sqs:GetQueueAttributes",
        "sqs:ChangeMessageVisibility",
      ]
      Resource = [
        aws_sqs_queue.transcode.arn,
        aws_sqs_queue.transcode_dlq.arn,
      ]
    }]
  })
}
