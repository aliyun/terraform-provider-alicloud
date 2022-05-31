data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_oos_application" "default" {
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  application_name  = "example_value"
  description       = "example_value"
  tags = {
    Created = "TF"
  }
}

resource "alicloud_oos_application_group" "default" {
  application_group_name = var.name
  application_name       = alicloud_oos_application.default.id
  deploy_region_id       = "example_value"
  description            = "example_value"
  import_tag_key         = "example_value"
  import_tag_value       = "example_value"
}

