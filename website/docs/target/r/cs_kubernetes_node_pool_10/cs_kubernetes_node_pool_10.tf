resource "alicloud_cs_kubernetes_node_pool" "default" {
  name                 = "existing-node"
  cluster_id           = alicloud_cs_managed_kubernetes.default.0.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  instance_charge_type = "PostPaid"

  # add existing node to nodepool
  instances = ["instance_id_01", "instance_id_02", "instance_id_03"]
  # default is false
  format_disk = false
  # default is true
  keep_instance_name = true
}
