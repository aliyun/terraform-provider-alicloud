provider "alicloud" {
  region = "cn-shanghai"
}
resource "alicloud_vpc" "vpc" {
  vpc_name   = "tf_test_foo"
  cidr_block = "172.16.0.0/12"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = "tf-testAccCsgName"
}
resource "alicloud_cloud_storage_gateway_storage_bundle" "example" {
  storage_bundle_name = "example_value"
}
resource "alicloud_cloud_storage_gateway_gateway" "default" {
  description              = "tf-acctestDesalone"
  gateway_class            = "Standard"
  type                     = "File"
  payment_type             = "PayAsYouGo"
  vswitch_id               = alicloud_vswitch.default.id
  release_after_expiration = false
  public_network_bandwidth = 40
  storage_bundle_id        = alicloud_cloud_storage_gateway_storage_bundle.example.id
  location                 = "Cloud"
  gateway_name             = "tf-acctestGatewayName"
}
