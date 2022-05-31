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
locals {
  vswitch_id = data.alicloud_vswitches.default.ids[0]
}

resource "alicloud_vpn_gateway" "default" {
  name                 = var.name
  vpc_id               = data.alicloud_vpcs.default.ids.0
  bandwidth            = 10
  enable_ssl           = true
  enable_ipsec         = true
  ssl_connections      = 5
  instance_charge_type = "PrePaid"
  vswitch_id           = local.vswitch_id
}

resource "alicloud_vpn_ipsec_server" "example" {
  client_ip_pool    = "example_value"
  ipsec_server_name = "example_value"
  local_subnet      = "example_value"
  vpn_gateway_id    = alicloud_vpn_gateway.default.id
}
