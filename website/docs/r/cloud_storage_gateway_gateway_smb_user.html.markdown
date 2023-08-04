---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_gateway_smb_user"
sidebar_current: "docs-alicloud-resource-cloud-storage-gateway-gateway-smb-user"
description: |-
  Provides a Alicloud Cloud Storage Gateway Gateway SMB User resource.
---

# alicloud_cloud_storage_gateway_gateway_smb_user

Provides a Cloud Storage Gateway Gateway SMB User resource.

For information about Cloud Storage Gateway Gateway SMB User and how to use it, see [What is Gateway SMB User](https://www.alibabacloud.com/help/en/cloud-storage-gateway/latest/creategatewaysmbuser).

-> **NOTE:** Available since v1.142.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_uuid" "default" {
}
resource "alicloud_cloud_storage_gateway_storage_bundle" "default" {
  storage_bundle_name = substr("tf-example-${replace(random_uuid.default.result, "-", "")}", 0, 16)
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_cloud_storage_gateway_gateway" "default" {
  gateway_name             = var.name
  description              = var.name
  gateway_class            = "Standard"
  type                     = "File"
  payment_type             = "PayAsYouGo"
  vswitch_id               = alicloud_vswitch.default.id
  release_after_expiration = false
  public_network_bandwidth = 40
  storage_bundle_id        = alicloud_cloud_storage_gateway_storage_bundle.default.id
  location                 = "Cloud"
}

resource "alicloud_cloud_storage_gateway_gateway_smb_user" "default" {
  username   = "example_username"
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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Gateway SMB User.
* `delete` - (Defaults to 1 mins) Used when delete the Gateway SMB User.

## Import

Cloud Storage Gateway Gateway SMB User can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_storage_gateway_gateway_smb_user.example <gateway_id>:<username>
```
