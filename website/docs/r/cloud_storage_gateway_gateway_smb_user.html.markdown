---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_gateway_smb_user"
sidebar_current: "docs-alicloud-resource-cloud-storage-gateway-gateway-smb-user"
description: |-
  Provides a Alicloud Cloud Storage Gateway Gateway SMB User resource.
---

# alicloud\_cloud\_storage\_gateway\_gateway\_smb\_user

Provides a Cloud Storage Gateway Gateway SMB User resource.

For information about Cloud Storage Gateway Gateway SMB User and how to use it, see [What is Gateway SMB User](https://www.alibabacloud.com/help/en/doc-detail/53972.htm).

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

```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, ForceNew) The Gateway ID of the Gateway SMB User.
* `password` - (Required, ForceNew) The password of the Gateway SMB User.
* `username` - (Required, ForceNew) The username of the Gateway SMB User.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Gateway SMB User. The value formats as `<gateway_id>:<username>`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Gateway SMB User.
* `delete` - (Defaults to 1 mins) Used when delete the Gateway SMB User.

## Import

Cloud Storage Gateway Gateway SMB User can be imported using the id, e.g.

```
$ terraform import alicloud_cloud_storage_gateway_gateway_smb_user.example <gateway_id>:<username>
```
