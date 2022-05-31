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

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_ecs_network_interface" "default" {
  network_interface_name = var.name
  vswitch_id             = data.alicloud_vswitches.default.ids.0
  security_group_ids     = [alicloud_security_group.default.id]
  description            = "Basic test"
  primary_ip_address     = cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 100)
  tags = {
    Created = "TF",
    For     = "Test",
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}
data "alicloud_account" "default" {}
resource "alicloud_ecs_network_interface_permission" "example" {
  account_id           = data.alicloud_account.default.id
  network_interface_id = alicloud_ecs_network_interface.default.id
  permission           = "InstanceAttach"
  force                = true
}
