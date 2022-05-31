resource "alicloud_cs_kubernetes_addon" "ack-node-problem-detector" {
  cluster_id = alicloud_cs_managed_kubernetes.default.0.id
  name       = "ack-node-problem-detector"
  version    = "1.2.7"
}

resource "alicloud_cs_kubernetes_addon" "nginx_ingress_controller" {
  cluster_id = var.cluster_id
  name       = "nginx-ingress-controller"
  version    = "v1.1.2-aliyun.2"
  // Specify custom configuration for addon. You can checkout the customizable configuration of the addon through data source alicloud_cs_kubernetes_addon_metadata.
  config = jsonencode(
    {
      CpuLimit              = ""
      CpuRequest            = "100m"
      EnableWebhook         = true
      HostNetwork           = false
      IngressSlbNetworkType = "internet"
      IngressSlbSpec        = "slb.s2.small"
      MemoryLimit           = ""
      MemoryRequest         = "200Mi"
      NodeSelector          = []
    }
  )
}
