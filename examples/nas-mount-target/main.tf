data "alicloud_nas_zones" "default" {
  file_system_type = "extreme"
}

locals {
  count_size = length(data.alicloud_nas_zones.default.zones)
  zone_id = data.alicloud_nas_zones.default.zones[local.count_size - 1].zone_id
}

resource "alicloud_vpc" "main" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "main" {
  vswitch_name = alicloud_vpc.main.vpc_name
  cidr_block   = alicloud_vpc.main.cidr_block
  vpc_id       = alicloud_vpc.main.id
  zone_id      = local.zone_id
}

resource "alicloud_nas_file_system" "main" {
  protocol_type     = "NFS"
  storage_type      = "advance"
  file_system_type  = "extreme"
  zone_id           = local.zone_id
  capacity          = "100"
}

resource "alicloud_nas_access_group" "main" {
  name             = "tf-testAccNasConfigName"
  type             = "Vpc"
  description      = "tf-testAccNasConfig"
  file_system_type = "extreme"
}

resource "alicloud_nas_mount_target" "main" {
  file_system_id    = alicloud_nas_file_system.main.id
  access_group_name = alicloud_nas_access_group.main.name
  vpc_id            = alicloud_vpc.main.id
  vswitch_id        = alicloud_vswitch.main.id
  network_type      = "Vpc"

}

