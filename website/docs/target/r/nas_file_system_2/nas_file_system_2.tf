data "alicloud_nas_zones" "default" {
  file_system_type = "cpfs"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_nas_zones.default.zones.0.zone_id
}

resource "alicloud_nas_file_system" "foo" {
  protocol_type    = "cpfs"
  storage_type     = "advance_200"
  file_system_type = "cpfs"
  capacity         = 3600
  description      = "tf-testacc"
  zone_id          = data.alicloud_nas_zones.default.zones.0.zone_id
  vpc_id           = data.alicloud_vpcs.default.ids.0
  vswitch_id       = data.alicloud_vswitches.default.ids.0
}
