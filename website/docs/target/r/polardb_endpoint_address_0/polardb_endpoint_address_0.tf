variable "creation" {
  default = "PolarDB"
}

variable "name" {
  default = "polardbconnectionbasic"
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

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "8.0"
  pay_type      = "PostPaid"
  db_node_class = "polar.mysql.x4.large"
  vswitch_id    = alicloud_vswitch.default.id
  description   = var.name
}

data "alicloud_polardb_endpoints" "default" {
  db_cluster_id = alicloud_polardb_cluster.default.id
}

resource "alicloud_polardb_endpoint_address" "endpoint" {
  db_cluster_id     = alicloud_polardb_cluster.default.id
  db_endpoint_id    = data.alicloud_polardb_endpoints.default.endpoints[0].db_endpoint_id
  connection_prefix = "testpolardbconn"
  net_type          = "Public"
}
