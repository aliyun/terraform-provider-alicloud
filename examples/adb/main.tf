// Zones data source for availability_zone
data "alicloud_zones" "default" {
  available_resource_creation = var.creation
}

resource "alicloud_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  availability_zone = data.alicloud_zones.default.zones.0.id
  name              = var.name
}
resource "alicloud_adb_cluster" "default" {
  db_cluster_version      = var.db_cluster_version
  db_cluster_category     = var.db_cluster_category
  db_cluster_network_type = "VPC"
  db_node_class           = var.db_node_class
  db_node_count           = var.db_node_count
  db_node_storage         = var.db_node_storage
  pay_type                = "PostPaid"
  description             = var.name
  vswitch_id              = alicloud_vswitch.default.id
}
