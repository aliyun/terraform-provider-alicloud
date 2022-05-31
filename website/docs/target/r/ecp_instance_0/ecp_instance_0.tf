data "alicloud_ecp_zones" "default" {
}

data "alicloud_ecp_instance_types" "default" {
}

locals {
  count_size               = length(data.alicloud_ecp_zones.default.zones)
  zone_id                  = data.alicloud_ecp_zones.default.zones[local.count_size - 1].zone_id
  instance_type_count_size = length(data.alicloud_ecp_instance_types.default.instance_types)
  instance_type            = data.alicloud_ecp_instance_types.default.instance_types[local.instance_type_count_size - 1].instance_type
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = local.zone_id
}

resource "alicloud_security_group" "group" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_ecp_key_pair" "default" {
  key_pair_name   = var.name
  public_key_body = "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg"
}

resource "alicloud_ecp_instance" "default" {
  instance_name     = var.name
  description       = var.name
  key_pair_name     = "${alicloud_ecp_key_pair.default.key_pair_name}"
  security_group_id = "${alicloud_security_group.group.id}"
  vswitch_id        = "${data.alicloud_vswitches.default.ids.0}"
  image_id          = "android_9_0_0_release_2851157_20211201.vhd"
  instance_type     = "${data.alicloud_ecp_instance_types.default.instance_types[local.instance_type_count_size - 1].instance_type}"
  vnc_password      = "Cp1234"
  charge_type       = "PayAsYouGo"
}
