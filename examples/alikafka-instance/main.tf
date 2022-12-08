resource "alicloud_alikafka_instance" "default" {
  name        = var.name
  partition_num = var.partition_num
  disk_type   = var.disk_type
  disk_size   = var.disk_size
  deploy_type = var.deploy_type
  io_max      = var.io_max
  eip_max     = var.eip_max
  paid_type   = var.paid_type
  spec_type   = var.spec_type
  vswitch_id  = var.vswitch_id
}