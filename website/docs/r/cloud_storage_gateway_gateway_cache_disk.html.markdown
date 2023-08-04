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

resource "random_uuid" "default" {
}
resource "alicloud_cloud_storage_gateway_storage_bundle" "default" {
  storage_bundle_name = substr("tf-example-${replace(random_uuid.default.result, "-", "")}", 0, 16)
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}
data "alicloud_cloud_storage_gateway_stocks" "default" {
  gateway_class = "Standard"
}
resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_cloud_storage_gateway_stocks.default.stocks.0.zone_id
  vswitch_name = var.name
}

resource "alicloud_cloud_storage_gateway_gateway" "default" {
  gateway_name             = var.name
  description              = var.name
  gateway_class            = "Standard"
  type                     = "File"
  payment_type             = "PayAsYouGo"
  vswitch_id               = alicloud_vswitch.default.id
  release_after_expiration = true
  public_network_bandwidth = 40
  storage_bundle_id        = alicloud_cloud_storage_gateway_storage_bundle.default.id
  location                 = "Cloud"
}

resource "alicloud_cloud_storage_gateway_gateway_cache_disk" "default" {
  cache_disk_category   = "cloud_efficiency"
  gateway_id            = alicloud_cloud_storage_gateway_gateway.default.id
  cache_disk_size_in_gb = 50
}
```

## Argument Reference

The following arguments are supported:

* `cache_disk_category` - (Optional, ForceNew) The cache disk type. Valid values: `cloud_efficiency`, `cloud_ssd`.
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

```shell
$ terraform import alicloud_cloud_storage_gateway_gateway_cache_disk.example <gateway_id>:<cache_id>:<local_file_path>
```