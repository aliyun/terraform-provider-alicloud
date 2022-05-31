variable "instance_name" {
  default = "alikafkaInstanceName"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = "172.16.0.0/24"
  zone_id    = data.alicloud_zones.default.zones[0].id
}

resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_alikafka_instance" "default" {
  name           = var.instance_name
  topic_quota    = "50"
  disk_type      = "1"
  disk_size      = "500"
  deploy_type    = "4"
  io_max         = "20"
  vswitch_id     = alicloud_vswitch.default.id
  security_group = alicloud_security_group.default.id
}
