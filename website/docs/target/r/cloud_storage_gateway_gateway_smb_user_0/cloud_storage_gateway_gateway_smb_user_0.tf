data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_cloud_storage_gateway_storage_bundle" "example" {
  storage_bundle_name = "example_value"
}

resource "alicloud_cloud_storage_gateway_gateway" "default" {
  description              = "tf-acctestDesalone"
  gateway_class            = "Standard"
  type                     = "File"
  payment_type             = "PayAsYouGo"
  vswitch_id               = data.alicloud_vswitches.default.ids.0
  release_after_expiration = false
  public_network_bandwidth = 40
  storage_bundle_id        = alicloud_cloud_storage_gateway_storage_bundle.example.id
  location                 = "Cloud"
  gateway_name             = "example_value"
}

resource "alicloud_cloud_storage_gateway_gateway_smb_user" "default" {
  username   = "your_username"
  password   = "password"
  gateway_id = alicloud_cloud_storage_gateway_gateway.default.id
}

