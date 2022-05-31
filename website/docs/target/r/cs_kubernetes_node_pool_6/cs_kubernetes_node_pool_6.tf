resource "alicloud_cs_kubernetes_node_pool" "default" {
  name                 = var.name
  cluster_id           = alicloud_cs_managed_kubernetes.default.0.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alicloud_key_pair.default.key_name
  # use PrePaid
  instance_charge_type = "PrePaid"
  period               = 1
  period_unit          = "Month"
  auto_renew           = true
  auto_renew_period    = 1

  # open cloud monitor
  install_cloud_monitor = true

  # enable auto-scaling
  scaling_config {
    min_size = 1
    max_size = 10
    type     = "cpu"
  }
}
