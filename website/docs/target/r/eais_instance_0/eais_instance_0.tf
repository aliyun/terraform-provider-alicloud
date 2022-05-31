variable "name" {
  default = "%v"
}
data "alicloud_vpcs" "default" {
  cidr_block = "172.16.0.0/12"
}
resource "alicloud_vpc" "default" {
  count      = length(data.alicloud_vpcs.default.ids) > 0 ? 0 : 1
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}
data "alicloud_vswitches" "default" {
  vpc_id  = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
  zone_id = "cn-hangzhou-h"
}
resource "alicloud_vswitch" "default" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  zone_id      = "cn-hangzhou-h"
  vswitch_name = var.name
}
resource "alicloud_security_group" "default" {
  name        = var.name
  description = "tf test"
  vpc_id      = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
}
resource "alicloud_eais_instance" "default" {
  instance_type     = "eais.ei-a6.4xlarge"
  instance_name     = var.name
  security_group_id = alicloud_security_group.default.id
  vswitch_id        = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : alicloud_vswitch.default[0].id
}
