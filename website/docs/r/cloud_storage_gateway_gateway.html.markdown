---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_gateway"
sidebar_current: "docs-alicloud-resource-cloud-storage-gateway-gateway"
description: |-
  Provides a Alicloud Cloud Storage Gateway Gateway resource.
---

# alicloud_cloud_storage_gateway_gateway

Provides a Cloud Storage Gateway Gateway resource.

For information about Cloud Storage Gateway Gateway and how to use it, see [What is Gateway](https://www.alibabacloud.com/help/en/csg/developer-reference/api-mnz46x).

-> **NOTE:** Available since v1.132.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_storage_gateway_gateway&exampleId=c72b5954-f4f8-33f9-c4bc-b4016d22ceebc8beb4f1&activeTab=example&spm=docs.r.cloud_storage_gateway_gateway.0.c72b5954f4&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_cloud_storage_gateway_storage_bundle" "default" {
  storage_bundle_name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "${var.name}-${random_integer.default.result}"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = "${var.name}-${random_integer.default.result}"
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.192.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_cloud_storage_gateway_gateway" "default" {
  storage_bundle_id        = alicloud_cloud_storage_gateway_storage_bundle.default.id
  type                     = "File"
  location                 = "Cloud"
  gateway_name             = var.name
  gateway_class            = "Standard"
  vswitch_id               = alicloud_vswitch.default.id
  public_network_bandwidth = 50
  payment_type             = "PayAsYouGo"
  description              = var.name
}
```

## Argument Reference

The following arguments are supported:

* `storage_bundle_id` - (Required, ForceNew) The ID of the gateway cluster.
* `type` - (Required, ForceNew) The type of the gateway. Valid values: `File`, `Iscsi`.
* `location` - (Required, ForceNew) The location of the gateway. Valid values: `Cloud`, `On_Premise`.
* `gateway_name` - (Required) The name of the gateway. The name must be `1` to `60` characters in length and can contain letters, digits, periods (.), underscores (_), and hyphens (-). It must start with a letter.
* `gateway_class` - (Optional) The specification of the gateway. Valid values: `Basic`, `Standard`, `Enhanced`, `Advanced`. **NOTE:** If `location` is set to `Cloud`, `gateway_class` is required. Otherwise, `gateway_class` will be ignored. If `payment_type` is set to `Subscription`, `gateway_class` cannot be modified.
* `vswitch_id` - (Optional, ForceNew) The ID of the VSwitch. **NOTE:** If `location` is set to `Cloud`, `vswitch_id` is required. Otherwise, `vswitch_id` will be ignored.
* `public_network_bandwidth` - (Optional, Int) The public bandwidth of the gateway. Default value: `5`. Valid values: `5` to `200`. **NOTE:** `public_network_bandwidth` is only valid when `location` is `Cloud`. If `payment_type` is set to `Subscription`, `public_network_bandwidth` cannot be modified.
* `payment_type` - (Optional, ForceNew) The Payment type of gateway. Valid values: `PayAsYouGo`, `Subscription`. **NOTE:** From version 1.233.0, `payment_type` can be set to `Subscription`.
* `description` - (Optional) The description of the gateway.
* `release_after_expiration` - (Optional, Bool, ForceNew) Specifies whether to release the gateway after the subscription expires. Valid values:
  - `true`: The gateway is released after it expires. The gateway becomes unavailable seven days after the gateway expires.
  - `false`: The gateway is not released after it expires. After the gateway expires, its billing method is changed to `PayAsYouGo`.
-> **NOTE:** `release_after_expiration` is only valid when `payment_type` is `Subscription`.
* `reason_type` - (Optional) The type of the reason why you want to delete the gateway.
* `reason_detail` - (Optional) The detailed reason why you want to delete the gateway.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Gateway.
* `status` - The status of the Gateway.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Gateway.
* `update` - (Defaults to 5 mins) Used when update the Gateway.
* `delete` - (Defaults to 5 mins) Used when delete the Gateway.

## Import

Cloud Storage Gateway Gateway can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_storage_gateway_gateway.example <id>
```
