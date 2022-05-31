data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
resource "alicloud_vpc" "default" {
  count    = length(data.alicloud_vpcs.default.ids) > 0 ? 0 : 1
  vpc_name = "example_value"
}
data "alicloud_vswitches" "default" {
  vpc_id = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
}
resource "alicloud_vswitch" "default" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = "example_value"
}
resource "alicloud_service_mesh_service_mesh" "example" {
  service_mesh_name = "example_value"
  network {
    vpc_id        = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
    vswitche_list = [length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : alicloud_vswitch.default[0].id]
  }
}

