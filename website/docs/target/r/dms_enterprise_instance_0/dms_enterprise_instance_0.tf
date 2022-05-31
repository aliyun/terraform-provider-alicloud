resource "alicloud_dms_enterprise_instance" "default" {
  tid               = "12345"
  instance_type     = "MySQL"
  instance_source   = "RDS"
  network_type      = "VPC"
  env_type          = "test"
  host              = "rm-uf648hgsxxxxxx.mysql.rds.aliyuncs.com"
  port              = 3306
  database_user     = "your_user_name"
  database_password = "Yourpassword123"
  instance_name     = "your_alias_name"
  dba_uid           = "1182725234xxxxxxx"
  safe_rule         = "自由操作"
  query_timeout     = 60
  export_timeout    = 600
  ecs_region        = "cn-shanghai"
}
