// Zones data source for availability_zone
data "alicloud_zones" "default" {
  available_resource_creation = "Cassandra"
}

resource "alicloud_cassandra_cluster" "default" {
  cluster_name        = "tf-cassandra-cluster-example"
  zone_id             = var.availability_zone == "" ? data.alicloud_zones.default.zones[length(data.alicloud_zones.default.ids) - 1].id : var.availability_zone
  auto_renew          = var.auto_renew
  auto_renew_period   = var.auto_renew_period
  data_center_name    = var.dc_name_1
  disk_size           = var.disk_size
  disk_type           = var.disk_type
  instance_type       = var.instance_type
  major_version       = var.major_version
  node_count          = var.node_count
  password            = var.password
  pay_type            = var.pay_type
  vswitch_id          = var.vswitch_id
  maintain_start_time = var.maintain_start_time
  maintain_end_time   = var.maintain_end_time
  enable_public       = var.enable_public
  ip_white            = var.ip_white
}
