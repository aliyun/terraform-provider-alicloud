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
  default = "tf-testAccSagQosConfigName"
}

resource "alicloud_sag_qos" "default" {
  name      = var.name
  sag_count = "0"
}

resource "alicloud_sag_qos_car" "default" {
  qos_id            = "${alicloud_sag_qos.default.id}"
  name              = "tf-testSagQosCar"
  description       = "tf-testSagQosCar"
  priority          = "1"
  limit_type        = "Absolute"
  min_bandwidth_abs = "10"
  max_bandwidth_abs = "20"
}