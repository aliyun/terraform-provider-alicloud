variable "name" {
  default = "polardbClusterconfig"
}

variable "creation" {
  default = "PolarDB"
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

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "5.6"
  db_node_class = "polar.mysql.x4.medium"
  pay_type      = "PostPaid"
  description   = var.name
  vswitch_id    = alicloud_vswitch.default.id
  db_cluster_ip_array {
    db_cluster_ip_array_name = "default"
    security_ips             = ["1.2.3.4", "1.2.3.5"]
  }
  db_cluster_ip_array {
    db_cluster_ip_array_name = "test_ips1"
    security_ips             = ["1.2.3.6"]
  }
}
