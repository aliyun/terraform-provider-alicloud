resource "alicloud_edas_cluster" "default" {
  cluster_name      = var.cluster_name
  cluster_type      = var.cluster_type
  network_mode      = var.network_mode
  logical_region_id = var.logical_region_id
  vpc_id            = var.vpc_id
}