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

resource "alicloud_sag_dnat_entry" "default" {
  sag_id        = "sag-3rb1t3iagy3w0zgwy9"
  type          = "Intranet"
  ip_protocol   = "tcp"
  external_ip   = "1.0.0.2"
  external_port = "1"
  internal_ip   = "10.0.0.2"
  internal_port = "20"
}