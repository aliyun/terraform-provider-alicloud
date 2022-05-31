data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}

resource "alicloud_ecs_image_component" "example" {
  component_type       = "Build"
  content              = "RUN yum update -y"
  description          = "example_value"
  image_component_name = "example_value"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  system_type          = "Linux"
  tags = {
    Created = "TF"
  }
}
