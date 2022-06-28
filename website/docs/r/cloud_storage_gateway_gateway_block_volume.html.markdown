---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_gateway_block_volume"
sidebar_current: "docs-alicloud-resource-cloud-storage-gateway-gateway-block-volume"
description: |-
  Provides a Alicloud Cloud Storage Gateway Gateway Block Volume resource.
---

# alicloud\_cloud\_storage\_gateway\_gateway\_block\_volume

Provides a Cloud Storage Gateway Gateway Block Volume resource.

For information about Cloud Storage Gateway Gateway Block Volume and how to use it, see [What is Gateway Block Volume](https://www.alibabacloud.com/help/en/doc-detail/53972.htm).

-> **NOTE:** Available in v1.144.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tftest"
}

data "alicloud_cloud_storage_gateway_stocks" "default" {
  gateway_class = "Standard"
}

resource "alicloud_vpc" "vpc" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_cloud_storage_gateway_stocks.default.stocks.0.zone_id
  vswitch_name = var.name
}

resource "alicloud_cloud_storage_gateway_storage_bundle" "default" {
  storage_bundle_name = var.name
}

resource "alicloud_cloud_storage_gateway_gateway" "default" {
  description              = "tf-acctestDesalone"
  gateway_class            = "Standard"
  type                     = "Iscsi"
  payment_type             = "PayAsYouGo"
  vswitch_id               = alicloud_vswitch.default.id
  release_after_expiration = true
  public_network_bandwidth = 10
  storage_bundle_id        = alicloud_cloud_storage_gateway_storage_bundle.default.id
  location                 = "Cloud"
  gateway_name             = var.name
}


resource "alicloud_cloud_storage_gateway_gateway_cache_disk" "default" {
  cache_disk_category   = "cloud_efficiency"
  gateway_id            = alicloud_cloud_storage_gateway_gateway.default.id
  cache_disk_size_in_gb = 50
}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
  acl    = "public-read-write"
}

resource "alicloud_cloud_storage_gateway_gateway_block_volume" "default" {
  cache_mode                = "Cache"
  chap_enabled              = true
  chap_in_user              = var.name
  chap_in_password          = var.name
  chunk_size                = "8192"
  gateway_block_volume_name = var.name
  gateway_id                = alicloud_cloud_storage_gateway_gateway.default.id
  local_path                = alicloud_cloud_storage_gateway_gateway_cache_disk.default.local_path
  oss_bucket_name           = alicloud_oss_bucket.default.bucket
  oss_bucket_ssl            = true
  oss_endpoint              = alicloud_oss_bucket.default.extranet_endpoint
  protocol                  = "iSCSI"
  size                      = 100
}
```

## Argument Reference

The following arguments are supported:

* `cache_mode` - (Optional, Computed, ForceNew) The Block volume set mode to cache mode. Value values: `Cache`, `WriteThrough`.
* `chap_enabled` - (Optional, Computed) Whether to enable iSCSI access of CHAP authentication, which currently supports both CHAP inbound authentication.  Default value: `false`.
* `chap_in_password` - (Optional) The password for inbound authentication when the block volume enables iSCSI access to CHAP authentication. **NOTE:** When the `chap_enabled` is  `true` is,The `chap_in_password` is valid.
* `chap_in_user` - (Optional) The Inbound CHAP user. The `chap_in_user` must be 1 to 32 characters in length, and can contain letters and digits. **NOTE:** When the `chap_enabled` is  `true` is,The `chap_in_password` is valid. 
* `chunk_size` - (Optional, Computed, ForceNew) The Block volume storage allocation unit.  Valid values: `8192`, `16384`, `32768`, `65536`, `131072`. Default value: `32768`. Unit: `Byte`.
* `gateway_block_volume_name` - (Required, ForceNew) The Block volume name. The name must be 1 to 32 characters in length, and can contain lower case letters and digits.
* `gateway_id` - (Required, ForceNew) The Gateway ID.
* `is_source_deletion` - (Optional) Whether to delete the source data. Default value `true`. **NOTE:** When `is_source_deletion` is `true`, the data in the OSS Bucket on the cloud is also deleted when deleting the block gateway volume. Please operate with caution.
* `local_path` - (Optional, ForceNew) The Cache disk to local path. **NOTE:**  When the `cache_mode` is  `Cache` is,The `chap_in_password` is valid.
* `oss_bucket_name` - (Required, ForceNew) The name of the OSS Bucket. 
* `oss_bucket_ssl` - (Optional, Computed, ForceNew) Whether to enable SSL access your OSS Buckets. Default value: `true`.
* `oss_endpoint` - (Required, ForceNew) The endpoint of the OSS Bucket.
* `protocol` - (Required, ForceNew) The Protocol. Valid values: `iSCSI`.
* `recovery` - (Optional) The recovery.
* `size` - (Optional) The Volume size. Valid values: `1` to `262144`. Unit: `Byte`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Gateway Block Volume. The value formats as `<gateway_id>:<index_id>`.
* `index_id` - The ID of the index.
* `status` - The status of volume. Valid values: 
  - `0`: Normal condition.
  - `1`: Failed to create volume.
  - `2`: Failed to delete volume.
  - `3`: Failed to enable target.
  - `4`: Failed to disable target.
  - `5`: Database error.
  - `6`: Failed to enable cache.
  - `7`: Failed to disable cache.
  - `8`: System error.

## Import

Cloud Storage Gateway Gateway Block Volume can be imported using the id, e.g.

```
$ terraform import alicloud_cloud_storage_gateway_gateway_block_volume.example <gateway_id>:<index_id>
```