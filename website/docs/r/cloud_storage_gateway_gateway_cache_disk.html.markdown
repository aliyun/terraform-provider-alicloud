---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_gateway_cache_disk"
sidebar_current: "docs-alicloud-resource-cloud-storage-gateway-gateway-cache-disk"
description: |-
  Provides a Alicloud Cloud Storage Gateway Gateway Cache Disk resource.
---

# alicloud_cloud_storage_gateway_gateway_cache_disk

Provides a Cloud Storage Gateway Gateway Cache Disk resource.

For information about Cloud Storage Gateway Gateway Cache Disk and how to use it, see [What is Gateway Cache Disk](https://www.alibabacloud.com/help/en/cloud-storage-gateway/latest/creategatewaycachedisk).

-> **NOTE:** Available since v1.144.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_cloud_storage_gateway_stocks" "default" {
  gateway_class = "Standard"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_cloud_storage_gateway_stocks.default.stocks.0.zone_id
  vswitch_name = var.name
}

resource "alicloud_cloud_storage_gateway_storage_bundle" "default" {
  storage_bundle_name = var.name
}

resource "alicloud_cloud_storage_gateway_gateway" "default" {
  description              = var.name
  gateway_class            = "Standard"
  type                     = "File"
  payment_type             = "PayAsYouGo"
  vswitch_id               = alicloud_vswitch.default.id
  release_after_expiration = true
  storage_bundle_id        = alicloud_cloud_storage_gateway_storage_bundle.default.id
  location                 = "Cloud"
  gateway_name             = var.name
}

resource "alicloud_cloud_storage_gateway_gateway_cache_disk" "default" {
  gateway_id            = alicloud_cloud_storage_gateway_gateway.default.id
  cache_disk_size_in_gb = 50
  cache_disk_category   = "cloud_efficiency"
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, ForceNew) The ID of the gateway.
* `cache_disk_size_in_gb` - (Required, Int) The capacity of the cache disk.
* `cache_disk_category` - (Optional, ForceNew) The type of the cache disk. Valid values: `cloud_efficiency`, `cloud_ssd`, `cloud_essd`. **NOTE:** From version 1.226.1, `cache_disk_category` can be set to `cloud_essd`.
* `performance_level` - (Optional, ForceNew, Available since v1.226.1) The performance level (PL) of the Enterprise SSD (ESSD). Valid values: `PL1`, `PL2`, `PL3`. **NOTE:** If `cache_disk_category` is set to `cloud_essd`, `performance_level` is required.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Gateway Cache Disk. It formats as `<gateway_id>:<cache_id>:<local_file_path>`.
* `cache_id` - The ID of the cache disk.
* `local_file_path` - The path of the cache disk.
* `status` - The status of the Gateway Cache Disk.

## Import

Cloud Storage Gateway Gateway Cache Disk can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_storage_gateway_gateway_cache_disk.example <gateway_id>:<cache_id>:<local_file_path>
```
