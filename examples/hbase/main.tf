// Zones data source for availability_zone
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

resource "alicloud_hbase_instance" "default" {
  name = "tf-hbase-instance-test4"
  zone_id = var.availability_zone == "" ? data.alicloud_zones.default.zones[0].id : var.availability_zone
  engine_version = var.engine_version
  master_instance_type = var.master_instance_type
  core_instance_type = var.core_instance_type
  core_instance_quantity = var.core_instance_quantity
  core_disk_type = var.core_disk_type
  pay_type = var.pay_type
  duration = var.duration
  auto_renew = var.auto_renew
  vswitch_id = var.vswitch_id
  is_cold_storage = var.is_cold_storage
}
