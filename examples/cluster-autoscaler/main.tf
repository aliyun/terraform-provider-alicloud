resource "alicloud_cs_kubernetes_autoscaler" "default" {
  cluster_id              = "${var.cluster_id}"
  nodepools {
        id                = "scaling_group_id"
        taints            = "c=d:NoSchedule"
        labels            = "a=b"
  }
  utilization             = "${var.utilization}"
  cool_down_duration      = "${var.cool_down_duration}"
  defer_scale_in_duration = "${var.defer_scale_in_duration}"
}