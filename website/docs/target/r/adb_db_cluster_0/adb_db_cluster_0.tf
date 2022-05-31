variable "name" {
  default = "adbClusterconfig"
}

variable "creation" {
  default = "ADB"
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

resource "alicloud_adb_db_cluster" "this" {
  db_cluster_category = "Cluster"
  db_cluster_class    = "C8"
  db_node_count       = "4"
  db_node_storage     = "400"
  mode                = "reserver"
  db_cluster_version  = "3.0"
  payment_type        = "PayAsYouGo"
  vswitch_id          = alicloud_vswitch.default.id
  description         = "Test new adb again."
  maintain_time       = "23:00Z-00:00Z"
  tags = {
    Created = "TF-update"
    For     = "acceptance-test-update"
  }
  resource_group_id = "rg-aek2s7ylxx6****"
  security_ips      = ["10.168.1.12", "10.168.1.11"]
}
