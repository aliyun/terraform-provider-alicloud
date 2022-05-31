resource "alicloud_vpc" "example" {
  vpc_name   = "vpc-123456"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vpc_id       = alicloud_vpc.example.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.example.zones.0.id
  vswitch_name = "vpc-123456"
}

resource "alicloud_db_instance" "default" {
  engine              = "MySQL"
  engine_version      = "5.6"
  db_instance_class   = "rds.mysql.t1.small"
  db_instance_storage = "10"
  vswitch_id          = alicloud_vswitch.example.id
}

resource "alicloud_db_instance" "example" {
  engine              = "MySQL"
  engine_version      = "5.6"
  db_instance_class   = "rds.mysql.t1.small"
  db_instance_storage = "10"
  parameters {
    name  = "innodb_large_prefix"
    value = "ON"
  }
  parameters {
    name  = "connect_timeout"
    value = "50"
  }
}
