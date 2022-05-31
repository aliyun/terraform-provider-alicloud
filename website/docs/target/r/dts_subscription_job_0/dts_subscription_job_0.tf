variable "name" {
  default = "dtsSubscriptionJob"
}

variable "creation" {
  default = "Rds"
}

data "alicloud_zones" "default" {
  available_resource_creation = var.creation
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_db_instance" "instance" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.s1.small"
  instance_storage = "10"
  vswitch_id       = alicloud_vswitch.default.id
  instance_name    = var.name
}

resource "alicloud_db_database" "db" {
  count       = 2
  instance_id = alicloud_db_instance.instance.id
  name        = "tfaccountpri_${count.index}"
  description = "from terraform"
}

resource "alicloud_db_account" "account" {
  instance_id = alicloud_db_instance.instance.id
  name        = "tftestprivilege"
  password    = "Test12345"
  description = "from terraform"
}

resource "alicloud_db_account_privilege" "privilege" {
  instance_id  = alicloud_db_instance.instance.id
  account_name = alicloud_db_account.account.name
  privilege    = "ReadWrite"
  db_names     = alicloud_db_database.db.*.name
}


data "alicloud_vpcs" "default1" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default_1" {
  vpc_id = data.alicloud_vpcs.default.ids[0]
}

resource "alicloud_dts_subscription_job" "default" {
  dts_job_name                       = var.name
  payment_type                       = "PostPaid"
  source_endpoint_engine_name        = "MySQL"
  source_endpoint_region             = "cn-hangzhou"
  source_endpoint_instance_type      = "RDS"
  source_endpoint_instance_id        = alicloud_db_instance.instance.id
  source_endpoint_database_name      = "tfaccountpri_0"
  source_endpoint_user_name          = "tftestprivilege"
  source_endpoint_password           = "Test12345"
  db_list                            = <<EOF
        {"dtstestdata": {"name": "tfaccountpri_0", "all": true}}
    EOF
  subscription_instance_network_type = "vpc"
  subscription_instance_vpc_id       = data.alicloud_vpcs.default1.ids[0]
  subscription_instance_vswitch_id   = data.alicloud_vswitches.default_1.ids[0]
  status                             = "Normal"
}
