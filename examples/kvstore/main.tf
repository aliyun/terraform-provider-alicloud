data "alicloud_zones" "default" {
  available_resource_creation = "KVStore"
}

// VPC Resource for Module
resource "alicloud_vpc" "vpc" {
  count      = var.vpc_id == "" ? 1 : 0
  name       = var.vpc_name
  cidr_block = var.vpc_cidr
}

// VSwitch Resource for Module
resource "alicloud_vswitch" "vswitch" {
  count             = var.vswitch_id == "" ? 1 : 0
  availability_zone = var.availability_zone == "" ? data.alicloud_zones.default.zones[0].id : var.availability_zone
  name              = var.vswitch_name
  cidr_block        = var.vswitch_cidr
  vpc_id            = var.vpc_id == "" ? alicloud_vpc.vpc[0].id : var.vpc_id
}

resource "alicloud_kvstore_instance" "myredis" {
  instance_class = var.instance_class
  instance_name  = var.instance_name
  password       = var.password
  vswitch_id     = var.vswitch_id == "" ? alicloud_vswitch.vswitch[0].id : var.vswitch_id
  security_ips   = ["1.1.1.1", "2.2.2.2", "3.3.3.3"]
  vpc_auth_mode  = "Close"
  engine_version = "4.0"

  //Refer to https://help.aliyun.com/document_detail/43885.html
  parameters {
    # {
    #   name = "cluster_compat_enable"
    #   value = "1"
    # },
    name = "maxmemory-policy"

    value = "volatile-ttl"
  }
}

resource "alicloud_kvstore_backup_policy" "redisbackup" {
  instance_id   = alicloud_kvstore_instance.myredis.id
  backup_time   = "03:00Z-04:00Z"
  backup_period = ["Monday", "Wednesday", "Friday"]
}

resource "alicloud_kvstore_account" "account" {
  instance_id = alicloud_kvstore_instance.myredis.id
  account_name        = "tftestnormal"
  account_password    = "Test12345"
}

