# RDS を置くサブネットの集合（private・マルチAZ）
resource "aws_db_subnet_group" "main" {
  name       = "my-video-app-db-subnet"
  subnet_ids = [aws_subnet.private_a.id, aws_subnet.private_c.id]
  tags       = { Name = "my-video-app-db-subnet" }
}

resource "aws_db_instance" "main" {
  identifier     = "my-video-app-db"
  engine         = "postgres"
  engine_version = "18.3"
  instance_class = "db.t4g.micro"
  db_name        = "myvideoapp"
  password       = var.db_password

  allocated_storage     = 20
  max_allocated_storage = 1000 # ストレージ自動スケーリングON（実物に合わせる）
  storage_type          = "gp2"
  storage_encrypted     = true

  username = "postgres"
  port     = 5432
  # db_name は未設定（実物がnull＝既定のpostgres DBのみ。あえて書かない）

  multi_az                   = false
  availability_zone          = "ap-northeast-1a"
  publicly_accessible        = false
  db_subnet_group_name       = aws_db_subnet_group.main.name
  parameter_group_name       = "default.postgres18"
  backup_retention_period    = 1
  auto_minor_version_upgrade = true
  deletion_protection        = false

  vpc_security_group_ids = [aws_security_group.rds.id]

  # 将来 destroy する時に最終スナップショットを取らない（学習用）
  skip_final_snapshot = true

  # パスワードはAWSから読めない＝stateに入らない。差分扱いさせない。
  lifecycle {
    ignore_changes = [password]
  }

  # --- 監視まわり（実物に合わせる） ---
  copy_tags_to_snapshot = true

  monitoring_interval = 60
  monitoring_role_arn = "arn:aws:iam::104954589692:role/rds-monitoring-role"

  performance_insights_enabled          = true
  performance_insights_kms_key_id       = "arn:aws:kms:ap-northeast-1:104954589692:key/a32186bf-d9cb-477e-8f99-01c542780db1"
  performance_insights_retention_period = 7
}
