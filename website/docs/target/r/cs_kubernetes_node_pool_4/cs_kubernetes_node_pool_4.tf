resource "alicloud_cs_kubernetes_node_pool" "default" {
  name                 = var.name
  cluster_id           = alicloud_cs_managed_kubernetes.default.0.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alicloud_key_pair.default.key_name

  # automatic scaling node pool configuration.
  # With auto-scaling is enabled, the nodes in the node pool will be labeled with `k8s.aliyun.com=true` to prevent system pods such as coredns, metrics-servers from being scheduled to elastic nodes, and to prevent node shrinkage from causing business abnormalities.
  scaling_config {
    min_size = 1
    max_size = 10
  }

}
