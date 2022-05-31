resource "alicloud_cs_kubernetes_node_pool" "default" {
  name                 = "windows-np"
  cluster_id           = alicloud_cs_managed_kubernetes.default.0.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  instance_charge_type = "PostPaid"
  desired_size         = 1

  // if the instance platform is windows, the password is requered.
  password = "Hello1234"
  platform = "Windows"
  image_id = "${window_image_id}"
}
