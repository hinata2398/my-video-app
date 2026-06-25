# --- ALB: 公開HTTPを受ける ---
resource "aws_security_group" "alb" {
  name        = "my-video-app-alb"
  description = "ALB: public HTTP in"
  vpc_id      = aws_vpc.main.id
  tags        = { Name = "my-video-app-alb" }
}

resource "aws_vpc_security_group_ingress_rule" "alb_http" {
  security_group_id = aws_security_group.alb.id
  description       = "HTTP from anywhere"
  ip_protocol       = "tcp"
  from_port         = 80
  to_port           = 80
  cidr_ipv4         = "0.0.0.0/0"
}

resource "aws_vpc_security_group_egress_rule" "alb_all_out" {
  security_group_id = aws_security_group.alb.id
  ip_protocol       = "-1" # 全プロトコル
  cidr_ipv4         = "0.0.0.0/0"
}

# --- EC2(アプリ): ALBからの8080だけ受ける ---
resource "aws_security_group" "ec2" {
  name        = "my-video-app-ec2"
  description = "App: 8080 from ALB only"
  vpc_id      = aws_vpc.main.id
  tags        = { Name = "my-video-app-ec2" }
}

resource "aws_vpc_security_group_ingress_rule" "ec2_from_alb" {
  security_group_id            = aws_security_group.ec2.id
  description                  = "App 8080 from ALB SG"
  ip_protocol                  = "tcp"
  from_port                    = 8080
  to_port                      = 8080
  referenced_security_group_id = aws_security_group.alb.id # ← SG-to-SG
}

resource "aws_vpc_security_group_egress_rule" "ec2_all_out" {
  security_group_id = aws_security_group.ec2.id
  ip_protocol       = "-1"
  cidr_ipv4         = "0.0.0.0/0"
}

# --- RDS: EC2(アプリ)からの5432だけ受ける ---
resource "aws_security_group" "rds" {
  name        = "my-video-app-rds"
  description = "RDS: 5432 from EC2 app SG only"
  vpc_id      = aws_vpc.main.id
  tags        = { Name = "my-video-app-rds" }
}

resource "aws_vpc_security_group_ingress_rule" "rds_from_ec2" {
  security_group_id            = aws_security_group.rds.id
  description                  = "PostgreSQL 5432 from EC2 app SG"
  ip_protocol                  = "tcp"
  from_port                    = 5432
  to_port                      = 5432
  referenced_security_group_id = aws_security_group.ec2.id # ← SG-to-SG
}

resource "aws_vpc_security_group_ingress_rule" "alb_https" {
  security_group_id = aws_security_group.alb.id
  description       = "HTTPS from anywhere"
  ip_protocol       = "tcp"
  from_port         = 443
  to_port           = 443
  cidr_ipv4         = "0.0.0.0/0"
}
