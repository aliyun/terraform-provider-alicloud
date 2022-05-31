data "alicloud_instance_types" "default" {
  instance_type_family = "ecs.g7"
}

data "alicloud_zones" "default" {
  available_resource_creation = "Instance"
  available_instance_type     = data.alicloud_instance_types.default.instance_types.0.id
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}

locals {
  vswitch_id = data.alicloud_vswitches.default.ids[0]
}

resource "alicloud_security_group" "default" {
  name        = "tf-test"
  description = "New security group"
  vpc_id      = data.alicloud_vpcs.default.ids.0
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

resource "alicloud_instance" "default" {
  count                = 2
  availability_zone    = data.alicloud_zones.default.zones.0.id
  instance_name        = "example_value"
  host_name            = "tf-testAcc"
  image_id             = data.alicloud_images.default.images.0.id
  instance_type        = data.alicloud_instance_types.default.instance_types.0.id
  security_groups      = [alicloud_security_group.default.id]
  vswitch_id           = local.vswitch_id
  system_disk_category = "cloud_essd"
}

resource "alicloud_ecs_network_interface" "default" {
  count                  = 2
  network_interface_name = "example_value"
  vswitch_id             = local.vswitch_id
  security_group_ids     = [alicloud_security_group.default.id]
}

resource "alicloud_ecs_network_interface_attachment" "default" {
  count                = 2
  instance_id          = element(alicloud_instance.default.*.id, count.index)
  network_interface_id = element(alicloud_ecs_network_interface.default.*.id, count.index)
}

resource "alicloud_vpc_traffic_mirror_filter" "default" {
  traffic_mirror_filter_name        = "example_value"
  traffic_mirror_filter_description = "example_value"
}

resource "alicloud_vpc_traffic_mirror_session" "default" {
  priority                           = 1
  virtual_network_id                 = 10
  traffic_mirror_session_description = "example_value"
  traffic_mirror_session_name        = "example_value"
  traffic_mirror_target_id           = alicloud_ecs_network_interface_attachment.default[0].network_interface_id
  traffic_mirror_source_ids          = [alicloud_ecs_network_interface_attachment.default[1].network_interface_id]
  traffic_mirror_filter_id           = alicloud_vpc_traffic_mirror_filter.default.id
  traffic_mirror_target_type         = "NetworkInterface"
}
