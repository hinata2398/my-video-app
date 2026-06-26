resource "aws_efs_file_system" "scratch" {
  lifecycle_policy {
    transition_to_ia = "AFTER_30_DAYS"
  }
}

resource "aws_efs_mount_target" "private_a" {
  file_system_id  = aws_efs_file_system.scratch.id
  subnet_id       = aws_subnet.private_a.id
  security_groups = [aws_security_group.efs.id]
}

resource "aws_efs_mount_target" "private_c" {
  file_system_id  = aws_efs_file_system.scratch.id
  subnet_id       = aws_subnet.private_c.id
  security_groups = [aws_security_group.efs.id]
}

resource "aws_security_group" "efs" {
  name   = "my-video-app-efs"
  vpc_id = aws_vpc.main.id
}

# EC2からNFS(2049)のみ許可
resource "aws_vpc_security_group_ingress_rule" "efs_nfs" {
  security_group_id            = aws_security_group.efs.id
  from_port                    = 2049
  to_port                      = 2049
  ip_protocol                  = "tcp"
  referenced_security_group_id = aws_security_group.ec2.id
}
