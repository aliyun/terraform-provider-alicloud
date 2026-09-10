resource "alicloud_edas_application_deployment" "default" {
  app_id          = var.app_id
  group_id        = var.group_id
  package_version = var.package_version
  war_url         = var.war_url
}