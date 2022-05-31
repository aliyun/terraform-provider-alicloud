variable "name" {
  default = "auto_provisioning_group"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_auto_provisioning_group" "default" {
  launch_template_id            = alicloud_ecs_launch_template.template.id
  total_target_capacity         = "4"
  pay_as_you_go_target_capacity = "1"
  spot_target_capacity          = "2"
  launch_template_config {
    instance_type     = "ecs.n1.small"
    vswitch_id        = alicloud_vswitch.default.id
    weighted_capacity = "2"
    max_price         = "2"
  }
}

resource "alicloud_ecs_launch_template" "template" {
  name              = var.name
  image_id          = data.alicloud_images.default.images[0].id
  instance_type     = "ecs.n1.tiny"
  security_group_id = alicloud_security_group.default.id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}
