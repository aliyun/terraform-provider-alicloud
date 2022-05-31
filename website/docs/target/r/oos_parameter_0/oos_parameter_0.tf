data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_oos_parameter" "example" {
  parameter_name = "my-Parameter"
  type           = "String"
  value          = "example_value"
  description    = "example_value"
  tags = {
    Created = "TF"
    For     = "OosParameter"
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
}

