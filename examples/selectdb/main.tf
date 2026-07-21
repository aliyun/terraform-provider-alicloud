// Zones data source for availability_zone
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
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}
resource "alicloud_selectdb_db_instance" "default" {
  db_instance_class       = "selectdb.xlarge"
  db_instance_description = var.name
  cache_size              = 200
  payment_type            = "PayAsYouGo"
  engine_minor_version    = "3.0.12"
  vpc_id                  = alicloud_vpc.default.id
  zone_id                 = data.alicloud_zones.default.zones.0.id
  vswitch_id              = alicloud_vswitches.default.id
}
resource "alicloud_selectdb_db_cluster" "default" {
  db_instance_id         = alicloud_selectdb_db_instance.default.id
  db_cluster_description = var.name
  db_cluster_class       = "selectdb.2xlarge"
  cache_size             = 400
  payment_type           = "PayAsYouGo"
}

data "alicloud_selectdb_db_instances" "default" {
  ids = [alicloud_selectdb_db_instance.default.id]
}

data "alicloud_selectdb_db_clusters" "default" {
  ids = [alicloud_selectdb_db_cluster.default.id]
}