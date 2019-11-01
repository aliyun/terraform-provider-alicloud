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

resource "alicloud_sag_qos_policy" "default" {
  qos_id            = "${alicloud_sag_qos.default.id}"
  name              = "tf-testSagQosPolicy"
  description       = "tf-testSagQosPolicy"
  priority          = "1"
  ip_protocol       = "ALL"
  source_cidr       = "10.10.10.0/24"
  source_port_range = "-1/-1"
  dest_cidr         = "192.168.10.0/24"
  dest_port_range   = "-1/-1"
  start_time        = "2019-10-27T16:41:33+0800"
  end_time          = "2019-10-28T16:41:33+0800"
}