data "alicloud_nas_zones" "default" {
  file_system_type = "cpfs"
}

locals {
  count_size = length(data.alicloud_nas_zones.default.zones)
  zone_id    = data.alicloud_nas_zones.default.zones[local.count_size - 1].zone_id
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = local.zone_id
}

resource "alicloud_nas_file_system" "default" {
  protocol_type    = "cpfs"
  storage_type     = "advance_200"
  file_system_type = "cpfs"
  capacity         = 3600
  description      = "tf-testacc"
  zone_id          = local.zone_id
  vpc_id           = data.alicloud_vpcs.default.ids.0
  vswitch_id       = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_nas_mount_target" "default" {
  file_system_id = alicloud_nas_file_system.default.id
  vswitch_id     = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_oss_bucket" "default" {
  bucket = "example_value"
  acl    = "private"
  tags = {
    cpfs-dataflow = "true"
  }
}

resource "alicloud_nas_fileset" "default" {
  depends_on       = ["alicloud_nas_mount_target.default"]
  file_system_id   = alicloud_nas_file_system.default.id
  description      = "example_value"
  file_system_path = "/example_path/"
}

resource "alicloud_nas_data_flow" "default" {
  fset_id              = alicloud_nas_fileset.default.fileset_id
  description          = "example_value"
  file_system_id       = alicloud_nas_file_system.default.id
  source_security_type = "SSL"
  source_storage       = join("", ["oss://", alicloud_oss_bucket.default.bucket])
  throughput           = 600
}
