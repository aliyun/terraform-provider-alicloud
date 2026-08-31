// Zones data source for availability_zone
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

resource "alicloud_hbase_instance" "default" {
  name                   = "tf-hbase-instance-example"
  zone_id                = var.availability_zone == "" ? data.alicloud_zones.default.zones[0].id : var.availability_zone
  engine                 = var.engine
  engine_version         = var.engine_version
  master_instance_type   = var.master_instance_type
  core_instance_type     = var.core_instance_type
  core_instance_quantity = var.core_instance_quantity
  core_disk_type         = var.core_disk_type
  core_disk_size         = var.core_disk_size
  pay_type               = var.pay_type
  duration               = var.duration
  auto_renew             = var.auto_renew
  vswitch_id             = var.vswitch_id
  cold_storage_size      = var.cold_storage_size
  maintain_start_time    = var.maintain_start_time
  maintain_end_time      = var.maintain_end_time
  deletion_protection    = var.deletion_protection
  immediate_delete_flag  = var.immediate_delete_flag
  ip_white               = var.ip_white
}
