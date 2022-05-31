variable "creation" {
  default = "PolarDB"
}

variable "name" {
  default = "dbaccountprivilegebasic"
}

data "alicloud_zones" "default" {
  available_resource_creation = var.creation
}

resource "alicloud_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = "172.16.0.0/24"
  zone_id    = data.alicloud_zones.default.zones.0.id
  name       = var.name
}

resource "alicloud_polardb_cluster" "cluster" {
  db_type       = "MySQL"
  db_version    = "8.0"
  pay_type      = "PostPaid"
  db_node_class = "polar.mysql.x4.large"
  vswitch_id    = alicloud_vswitch.default.id
  description   = var.name
}

resource "alicloud_polardb_database" "db" {
  db_cluster_id = alicloud_polardb_cluster.cluster.id
  db_name       = "tftestdatabase"
}

resource "alicloud_polardb_account" "account" {
  db_cluster_id       = alicloud_polardb_cluster.cluster.id
  account_name        = "tftestnormal"
  account_password    = "Test12345"
  account_description = var.name
}

resource "alicloud_polardb_account_privilege" "privilege" {
  db_cluster_id     = alicloud_polardb_cluster.cluster.id
  account_name      = alicloud_polardb_account.account.account_name
  account_privilege = "ReadOnly"
  db_names          = [alicloud_polardb_database.db.db_name]
}
