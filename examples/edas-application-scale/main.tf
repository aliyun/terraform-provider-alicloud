resource "alicloud_edas_application_scale" "default" {
  app_id       = var.app_id
  deploy_group = var.deploy_group
  ecu_info     = var.ecu_info
  force_status = var.force_status
}