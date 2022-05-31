resource "alicloud_cs_autoscaling_config" "default" {
  cluster_id = alicloud_cs_managed_kubernetes.default.0.id
  // configure auto scaling
  cool_down_duration        = "10m"
  unneeded_duration         = "10m"
  utilization_threshold     = "0.5"
  gpu_utilization_threshold = "0.5"
  scan_interval             = "30s"
  scale_down_enabled        = true
  expander                  = "least-waste"
}
