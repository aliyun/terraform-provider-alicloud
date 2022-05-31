variable "name" {
  default = "example_value"
}
data "alicloud_zones" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "example_value"
}
data "alicloud_vswitches" "default" {
  zone_id = data.alicloud_zones.default.zones.0.id
  vpc_id  = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}
resource "alicloud_lindorm_instance" "default" {
  disk_category              = "cloud_efficiency"
  payment_type               = "PayAsYouGo"
  zone_id                    = data.alicloud_zones.default.zones.0.id
  vswitch_id                 = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
  instance_name              = var.name
  table_engine_specification = "lindorm.c.xlarge"
  table_engine_node_count    = "2"
  instance_storage           = "480"
}
