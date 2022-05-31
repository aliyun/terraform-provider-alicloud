data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}
resource "alicloud_mse_gateway" "example" {
  gateway_name      = "example_value"
  replica           = 2
  spec              = "MSE_GTW_2_4_200_c"
  vswitch_id        = data.alicloud_vswitches.default.ids.0
  backup_vswitch_id = data.alicloud_vswitches.default.ids.1
  vpc_id            = data.alicloud_vpcs.default.ids.0
}
