data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones[0].id
  eni_amount        = 2
}

data "alicloud_images" "image" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

variable "name" {
default = "tf-testAccCcnConfigName"
}
variable "description" {
default = "tf-testAccCcnConfigDescription"
}

resource "alicloud_cloud_connect_network" "default" {
name        = var.name
description = var.description
cidr_block = "192.168.0.0/24"
is_default = true
}