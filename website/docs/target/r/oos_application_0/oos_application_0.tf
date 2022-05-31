variable "name" {
  default = "tftestacc"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_oos_application" "default" {
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  application_name  = var.name
  description       = var.name
  tags = {
    Created = "TF"
  }
}

