resource "alicloud_cs_kubernetes_addon" "ack-node-problem-detector" {
  cluster_id = alicloud_cs_managed_kubernetes.default.0.id
  name       = "ack-node-problem-detector"
  version    = "1.2.8" # upgrade from 1.2.7 to 1.2.8
}
