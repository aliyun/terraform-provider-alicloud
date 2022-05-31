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
  db_version    = "8.0"
  db_node_class = "polar.mysql.x4.large"
  pay_type      = "PostPaid"
  description   = var.name
  vswitch_id    = alicloud_vswitch.default.id
}

resource "alicloud_polardb_backup_policy" "policy" {
  db_cluster_id                               = alicloud_polardb_cluster.default.id
  preferred_backup_period                     = ["Tuesday", "Wednesday"]
  preferred_backup_time                       = "10:00Z-11:00Z"
  backup_retention_policy_on_cluster_deletion = "NONE"
}
