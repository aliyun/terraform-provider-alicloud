variable "name" {
  default = "tf-testAcc"
}

data "alicloud_zones" default {
  available_resource_creation = "Instance"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
  eni_amount        = 3
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/24"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "192.168.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vpc_id       = alicloud_vpc.default.id
}

resource "alicloud_security_group" "default" {
  name        = var.name
  description = "New security group"
  vpc_id      = alicloud_vpc.default.id
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

resource "alicloud_instance" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
  instance_name     = var.name
  host_name         = "tf-testAcc"
  image_id          = data.alicloud_images.default.images.0.id
  instance_type     = data.alicloud_instance_types.default.instance_types.0.id
  security_groups   = [alicloud_security_group.default.id]
  vswitch_id        = alicloud_vswitch.default.id
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_ecs_network_interface" "default" {
  network_interface_name = var.name
  vswitch_id             = alicloud_vswitch.default.id
  security_group_ids     = [alicloud_security_group.default.id]
  description            = "Basic test"
  primary_ip_address     = "192.168.0.2"
  tags = {
    Created = "TF",
    For     = "Test",
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}

resource "alicloud_ecs_network_interface_attachment" "default" {
  network_interface_id = alicloud_ecs_network_interface.default.id
  instance_id          = alicloud_instance.default.id
}

