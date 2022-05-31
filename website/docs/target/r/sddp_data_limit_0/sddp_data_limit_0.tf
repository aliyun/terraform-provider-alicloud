variable "name" {
  default = "tfaccxinmutes"
}

variable "region" {
  default = "cn-hangzhou"
}

variable "password" {
  default = "Test12345"
}

variable "database_name" {
  default = "tftestdatabase"
}

data "alicloud_db_zones" "default" {}

data "alicloud_db_instance_classes" "default" {
  engine         = "MySQL"
  engine_version = "5.6"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids[0]
  zone_id = data.alicloud_db_zones.default.zones[0].id
}

resource "alicloud_db_instance" "default" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = data.alicloud_db_instance_classes.default.instance_classes[0].instance_class
  instance_storage = "10"
  vswitch_id       = data.alicloud_vswitches.default.ids[0]
  instance_name    = var.name
}


locals {
  parent_id = join(".", [alicloud_db_instance.default.id, var.database_name])
}

resource "alicloud_rds_account" "default" {
  db_instance_id   = alicloud_db_instance.default.id
  account_name     = var.database_name
  account_password = var.password
}

resource "alicloud_db_database" "default" {
  instance_id = alicloud_db_instance.default.id
  name        = var.database_name
}

resource "alicloud_db_account_privilege" "default" {
  instance_id  = alicloud_db_instance.default.id
  account_name = alicloud_rds_account.default.name
  privilege    = "ReadWrite"
  db_names     = [alicloud_db_database.default.name]
}

resource "alicloud_sddp_data_limit" "default" {
  audit_status      = 0
  engine_type       = "MySQL"
  parent_id         = local.parent_id
  resource_type     = "RDS"
  user_name         = var.database_name
  password          = var.password
  port              = 3306
  service_region_id = var.region
  depends_on        = [alicloud_db_account_privilege.default]
}
