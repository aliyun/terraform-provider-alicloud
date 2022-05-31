variable "creation" {
  default = "KVStore"
}

variable "name" {
  default = "kvstoreinstancevpc"
}

data "alicloud_zones" "default" {
  available_resource_creation = var.creation
}

resource "alicloud_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_kvstore_instance" "default" {
  instance_class = "redis.master.small.default"
  instance_name  = var.name
  vswitch_id     = alicloud_vswitch.default.id
  private_ip     = "172.16.0.10"
  security_ips   = ["10.0.0.1"]
  instance_type  = "Redis"
  engine_version = "4.0"
}

resource "alicloud_kvstore_account" "example" {
  account_name     = "tftestnormal"
  account_password = "YourPassword_123"
  instance_id      = alicloud_kvstore_instance.default.id
}
