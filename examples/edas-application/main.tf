resource "alicloud_edas_application" "default" {
  application_name  = var.application_name
  cluster_id        = var.cluster_id
  package_type      = var.package_type
  build_pack_id     = var.build_pack_id
  descriotion       = var.descriotion
  health_check_url  = var.health_check_url
  logical_region_id = var.logical_region_id
  component_ids     = var.component_ids
  ecu_info          = var.ecu_info
  group_id          = var.group_id
  package_version   = var.package_version
  war_url           = var.war_url
}