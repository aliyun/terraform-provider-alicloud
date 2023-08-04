---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_gateway"
sidebar_current: "docs-alicloud-resource-cloud-storage-gateway-gateway"
description: |-
  Provides a Alicloud Cloud Storage Gateway Gateway resource.
---

# alicloud_cloud_storage_gateway_gateway

Provides a Cloud Storage Gateway: Gateway resource.

For information about Cloud Storage Gateway Gateway and how to use it, see [What is Gateway](https://www.alibabacloud.com/help/en/cloud-storage-gateway/latest/deploygateway).

-> **NOTE:** Available since v1.132.0.

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
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional)  the description of gateway.
* `gateway_class` - (Optional) the gateway class. the valid values: `Basic`, `Standard`,`Enhanced`,`Advanced`
* `gateway_name` - (Required) the name of gateway.
* `location` - (Required, ForceNew) gateway location. the valid values: `Cloud`, `On_Premise`.
* `payment_type` - (Optional, ForceNew) The Payment type of gateway. The valid value: `PayAsYouGo`.
* `public_network_bandwidth` - (Optional) The public network bandwidth of gateway. Valid values between `5` and `200`. Defaults to `5`.
* `reason_detail` - (Optional) The reason detail of gateway.
* `reason_type` - (Optional) The reason type when user deletes the gateway.
* `release_after_expiration` - (Optional) Whether to release the gateway due to expiration.
* `storage_bundle_id` - (Required, ForceNew) storage bundle id.
* `type` - (Required, ForceNew) gateway type. the valid values: `Type`, `Iscsi`.
* `vswitch_id` - (Optional, ForceNew) The vswitch id of gateway.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Gateway.
* `status` - gateway status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Gateway.

## Import

Cloud Storage Gateway Gateway can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_storage_gateway_gateway.example <id>
```
