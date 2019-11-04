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

resource "alicloud_sag_snat_entry" "default" {
  sag_id     = "sag-3rb1t3iagy3w0zgwy9"
  cidr_block = "192.168.7.0/24"
  snat_ip    = "192.0.0.2"
}