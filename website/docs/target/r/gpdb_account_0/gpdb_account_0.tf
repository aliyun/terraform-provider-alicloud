variable "name" {
  default = "tftestacc"
}
data "alicloud_gpdb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_gpdb_zones.default.zones.2.id
}

resource "alicloud_vswitch" "default" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_gpdb_zones.default.zones.3.id
  vswitch_name = var.name
}

resource "alicloud_gpdb_elastic_instance" "default" {
  engine                  = "gpdb"
  engine_version          = "6.0"
  seg_storage_type        = "cloud_essd"
  seg_node_num            = 4
  storage_size            = 50
  instance_spec           = "2C16G"
  db_instance_description = "Created by terraform"
  instance_network_type   = "VPC"
  payment_type            = "PayAsYouGo"
  vswitch_id              = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.default.*.id, [""])[0]
}

resource "alicloud_gpdb_account" "default" {
  account_name        = var.name
  db_instance_id      = alicloud_gpdb_elastic_instance.default.id
  account_password    = "TFTest123"
  account_description = var.name
}

