---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_gateway_cache_disk"
sidebar_current: "docs-alicloud-resource-cloud-storage-gateway-gateway-cache-disk"
description: |-
  Provides a Alicloud Cloud Storage Gateway Gateway Cache Disk resource.
---

# alicloud\_cloud\_storage\_gateway\_gateway\_cache\_disk

Provides a Cloud Storage Gateway Gateway Cache Disk resource.

For information about Cloud Storage Gateway Gateway Cache Disk and how to use it, see [What is Gateway Cache Disk](https://www.alibabacloud.com/help/zh/doc-detail/170294.htm).

-> **NOTE:** Available in v1.144.0+.

## Example Usage

Basic Usage

```terraform
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

```

## Argument Reference

The following arguments are supported:

* `cache_disk_category` - (Optional, Computed, ForceNew) The cache disk type. Valid values: `cloud_efficiency`, `cloud_ssd`.
* `cache_disk_size_in_gb` - (Required) size of the cache disk. Unit: `GB`. The upper limit of the basic gateway cache disk is `1` TB (`1024` GB), that of the standard gateway is `2` TB (`2048` GB), and that of other gateway cache disks is `32` TB (`32768` GB). The lower limit for the file gateway cache disk capacity is `40` GB, and the lower limit for the block gateway cache disk capacity is `20` GB.
* `gateway_id` - (Required, ForceNew) The ID of the gateway.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Gateway Cache Disk. The value formats as `<gateway_id>:<cache_id>:<local_file_path>`.
* `local_file_path` - The cache disk inside the device name.
* `status` - The status of the resource. Valid values: `0`, `1`, `2`. `0`: Normal. `1`: Is about to expire. `2`: Has expired.
* `cache_id` - The ID of the cache.

## Import

Cloud Storage Gateway Gateway Cache Disk can be imported using the id, e.g.

```
$ terraform import alicloud_cloud_storage_gateway_gateway_cache_disk.example <gateway_id>:<cache_id>:<local_file_path>
```