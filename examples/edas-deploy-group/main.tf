resource "alicloud_edas_deploy_group" "default" {
  app_id     = var.app_id
  group_name = var.group_name
}