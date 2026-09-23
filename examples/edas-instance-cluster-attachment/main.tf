resource "alicloud_edas_instance_cluster_attachment" "default" {
  cluster_id   = var.cluster_id
  instance_ids = var.instance_ids
}