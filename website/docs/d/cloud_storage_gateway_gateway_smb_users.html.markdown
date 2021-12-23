---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_gateway_smb_users"
sidebar_current: "docs-alicloud-datasource-cloud-storage-gateway-gateway-smb-users"
description: |-
  Provides a list of Cloud Storage Gateway Gateway SMB Users to the user.
---

# alicloud\_cloud\_storage\_gateway\_gateway\_smb\_users

This data source provides the Cloud Storage Gateway Gateway SMB Users of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.142.0+.

## Example Usage

Basic Usage

```terraform
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
data "alicloud_cloud_storage_gateway_gateway_smb_users" "ids" {
  gateway_id = alicloud_cloud_storage_gateway_gateway.default.id
  ids        = [alicloud_cloud_storage_gateway_gateway_smb_user.default.id]
}
output "cloud_storage_gateway_gateway_smb_user_id_1" {
  value = data.alicloud_cloud_storage_gateway_gateway_smb_users.ids.users.0.id
}

```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, ForceNew) The Gateway ID.
* `ids` - (Optional, ForceNew, Computed)  A list of Gateway SMB User IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Gateway SMB username.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `users` - A list of Cloud Storage Gateway SMB Users. Each element contains the following attributes:
	* `gateway_id` - The Gateway ID.
	* `id` - The ID of the Gateway SMB User.
	* `username` - The username of the Gateway SMB User.
