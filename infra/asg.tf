# 最新の Amazon Linux 2023 (x86_64) を動的取得（手作業のAMI ID固定の代わり）
data "aws_ami" "al2023" {
  most_recent = true
  owners      = ["amazon"]
  filter {
    name   = "name"
    values = ["al2023-ami-2023.*-x86_64"]
  }
}

resource "aws_launch_template" "app" {
  name_prefix   = "my-video-app-"
  image_id      = data.aws_ami.al2023.id
  instance_type = "t3.micro"

  iam_instance_profile {
    name = "ec2-ssm-role" # 既存のインスタンスプロファイル（SSM+SSMパラメータ読取）
  }

  vpc_security_group_ids = [aws_security_group.ec2.id] # ← TFの新EC2-SG

  user_data = base64encode(file("${path.module}/user_data.sh"))

  tag_specifications {
    resource_type = "instance"
    tags          = { Name = "my-video-app-app" }
  }
}

resource "aws_autoscaling_group" "app" {
  name                = "my-video-app-asg"
  vpc_zone_identifier = [aws_subnet.private_a.id, aws_subnet.private_c.id]
  target_group_arns   = [aws_lb_target_group.app.arn]

  min_size         = 2
  max_size         = 2
  desired_capacity = 2

  health_check_type         = "ELB" # ALBのヘルスチェックで生死判定
  health_check_grace_period = 600   # 起動～build～docker run に時間がかかるので猶予を長めに

  launch_template {
    id      = aws_launch_template.app.id
    version = "$Latest"
  }

  tag {
    key                 = "Name"
    value               = "my-video-app-app"
    propagate_at_launch = true
  }
}
