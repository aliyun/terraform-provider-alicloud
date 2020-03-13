data "alicloud_zones" "default" {
  available_disk_category = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones[0].id
eni_amount = 2
}

data "alicloud_images" "image" {
  name_regex = "^ubuntu_18.*64"
  most_recent = true
  owners = "system"
}

variable "name" {
  default = "tf-testAccSagAclConfigName"
}

resource "alicloud_sag_acl" "default" {
  name = var.name
}

resource "alicloud_sag_acl_rule" "default" {
  acl_id = "${alicloud_sag_acl.default.id}"
  description = "tf-testSagAclRule"
  policy = "accept"
  ip_protocol= "ALL"
  direction = "in"
  source_cidr = "10.10.10.0/24"
  source_port_range =    "-1/-1"
  dest_cidr = "192.168.10.0/24"
  dest_port_range = "-1/-1"
  priority = "1"
}
