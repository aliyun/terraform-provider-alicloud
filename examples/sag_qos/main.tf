data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones[0].id
  eni_amount        = 2
}

data "alicloud_images" "image" {
  name_regex  = "^ubuntu_18.*_64"
  most_recent = true
  owners      = "system"
}

variable "name" {
  default = "tf-testAccSagQosConfigName"
}


resource "alicloud_sag_qos" "default" {
  name      = var.name
  sag_count = "0"
}