data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}
data "alicloud_zones" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}
data "alicloud_instance_types" "default" {
  image_id = data.alicloud_images.default.ids.0
}
resource "alicloud_ecs_image_pipeline" "default" {
  add_account                = ["example_value"]
  base_image                 = data.alicloud_images.default.ids.0
  base_image_type            = "IMAGE"
  build_content              = "RUN yum update -y"
  delete_instance_on_failure = false
  image_name                 = "example_value"
  name                       = "example_value"
  description                = "example_value"
  instance_type              = data.alicloud_instance_types.default.ids.0
  resource_group_id          = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  internet_max_bandwidth_out = 20
  system_disk_size           = 40
  to_region_id               = ["cn-qingdao", "cn-zhangjiakou"]
  vswitch_id                 = data.alicloud_vswitches.default.ids.0
  tags = {
    Created = "TF"
    For     = "Acceptance-test"
  }
}
