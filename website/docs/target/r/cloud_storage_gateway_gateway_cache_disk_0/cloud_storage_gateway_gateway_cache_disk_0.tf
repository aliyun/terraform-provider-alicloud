data "alicloud_cloud_storage_gateway_stocks" "example" {
  gateway_class = "Standard"
}

resource "alicloud_vpc" "vpc" {
  vpc_name   = "example_value"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "example" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_cloud_storage_gateway_stocks.example.stocks.0.zone_id
  vswitch_name = "example_value"
}

resource "alicloud_cloud_storage_gateway_storage_bundle" "example" {
  storage_bundle_name = "example_value"
}

resource "alicloud_cloud_storage_gateway_gateway" "example" {
  description              = "tf-acctestDesalone"
  gateway_class            = "Standard"
  type                     = "File"
  payment_type             = "PayAsYouGo"
  vswitch_id               = alicloud_vswitch.example.id
  release_after_expiration = true
  public_network_bandwidth = 10
  storage_bundle_id        = alicloud_cloud_storage_gateway_storage_bundle.example.id
  location                 = "Cloud"
  gateway_name             = "example_value"
}

resource "alicloud_cloud_storage_gateway_gateway_cache_disk" "example" {
  cache_disk_category   = "cloud_efficiency"
  gateway_id            = alicloud_cloud_storage_gateway_gateways.example.id
  cache_disk_size_in_gb = 50
}

