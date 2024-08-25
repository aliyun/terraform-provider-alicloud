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

For information about Cloud Storage Gateway Gateway and how to use it, see [What is Gateway](https://www.alibabacloud.com/help/en/cloud-storage-gateway/latest/deploygateway).

-> **NOTE:** Available since v1.132.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_cloud_storage_gateway_gateway&exampleId=6d198b36-5fc2-43f8-5311-a0ce025d1f507d503740&activeTab=example&spm=docs.r.cloud_storage_gateway_gateway.0.6d198b365f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

* `location` - (Required, ForceNew) The location of the gateway. Valid values: `Cloud`, `On_Premise`.
* `storage_bundle_id` - (Required, ForceNew) The ID of the gateway cluster.
* `type` - (Required, ForceNew) The type of the gateway. Valid values: `File`, `Iscsi`.
* `gateway_name` - (Required) The name of the gateway.
* `description` - (Optional) The description of the gateway.
* `gateway_class` - (Optional) The specification of the gateway. Valid values: `Basic`, `Standard`,`Enhanced`,`Advanced`.
* `payment_type` - (Optional, ForceNew) The Payment type of gateway. Valid values: `PayAsYouGo`.
* `public_network_bandwidth` - (Optional, Int) The public network bandwidth of gateway. Default value: `5`. Valid values: `5` to `200`.
* `reason_detail` - (Optional) The reason detail of gateway.
* `reason_type` - (Optional) The reason type when user deletes the gateway.
* `release_after_expiration` - (Optional, Bool) Whether to release the gateway due to expiration.
* `vswitch_id` - (Optional, ForceNew) The ID of the vSwitch.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Gateway.
* `status` - The status of the Gateway.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Gateway.

## Import

Cloud Storage Gateway Gateway can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_storage_gateway_gateway.example <id>
```
