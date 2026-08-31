---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_gateway_block_volume"
sidebar_current: "docs-alicloud-resource-cloud-storage-gateway-gateway-block-volume"
description: |-
  Provides a Alicloud Cloud Storage Gateway Gateway Block Volume resource.
---

# alicloud_cloud_storage_gateway_gateway_block_volume

Provides a Cloud Storage Gateway Gateway Block Volume resource.

For information about Cloud Storage Gateway Gateway Block Volume and how to use it, see [What is Gateway Block Volume](https://www.alibabacloud.com/help/en/cloud-storage-gateway/latest/creategatewayblockvolume).

-> **NOTE:** Available since v1.144.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_storage_gateway_gateway_block_volume&exampleId=618b2bca-021c-40f8-be14-548407d915147685e0be&activeTab=example&spm=docs.r.cloud_storage_gateway_gateway_block_volume.0.618b2bca02&intl_lang=EN_US" target="_blank">
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

resource "alicloud_oss_bucket" "default" {
  bucket = substr("tf-example-${replace(random_uuid.default.result, "-", "")}", 0, 16)
}

resource "alicloud_oss_bucket_acl" "default" {
  bucket = alicloud_oss_bucket.default.bucket
  acl    = "public-read-write"
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
  type                     = "Iscsi"
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

resource "alicloud_cloud_storage_gateway_gateway_block_volume" "default" {
  cache_mode                = "Cache"
  chap_enabled              = false
  chunk_size                = "8192"
  gateway_block_volume_name = "example"
  gateway_id                = alicloud_cloud_storage_gateway_gateway.default.id
  local_path                = alicloud_cloud_storage_gateway_gateway_cache_disk.default.local_file_path
  oss_bucket_name           = alicloud_oss_bucket.default.bucket
  oss_bucket_ssl            = true
  oss_endpoint              = alicloud_oss_bucket.default.extranet_endpoint
  protocol                  = "iSCSI"
  size                      = 100
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cloud_storage_gateway_gateway_block_volume&spm=docs.r.cloud_storage_gateway_gateway_block_volume.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `cache_mode` - (Optional, ForceNew) The Block volume set mode to cache mode. Valid values: `Cache`, `WriteThrough`.
* `chap_enabled` - (Optional) Whether to enable iSCSI access of CHAP authentication, which currently supports both CHAP inbound authentication.  Default value: `false`.
* `chap_in_password` - (Optional) The password for inbound authentication when the block volume enables iSCSI access to CHAP authentication. **NOTE:** When the `chap_enabled` is  `true` is,The `chap_in_password` is valid.
* `chap_in_user` - (Optional) The Inbound CHAP user. The `chap_in_user` must be 1 to 32 characters in length, and can contain letters and digits. **NOTE:** When the `chap_enabled` is  `true` is,The `chap_in_password` is valid. 
* `chunk_size` - (Optional, ForceNew) The Block volume storage allocation unit.  Valid values: `8192`, `16384`, `32768`, `65536`, `131072`. Default value: `32768`. Unit: `Byte`.
* `gateway_block_volume_name` - (Required, ForceNew) The Block volume name. The name must be 1 to 32 characters in length, and can contain lower case letters and digits.
* `gateway_id` - (Required, ForceNew) The Gateway ID.
* `is_source_deletion` - (Optional) Whether to delete the source data. Default value `true`. **NOTE:** When `is_source_deletion` is `true`, the data in the OSS Bucket on the cloud is also deleted when deleting the block gateway volume. Please operate with caution.
* `local_path` - (Optional, ForceNew) The Cache disk to local path. **NOTE:**  When the `cache_mode` is  `Cache` is,The `chap_in_password` is valid.
* `oss_bucket_name` - (Required, ForceNew) The name of the OSS Bucket. 
* `oss_bucket_ssl` - (Optional, ForceNew) Whether to enable SSL access your OSS Buckets. Default value: `true`.
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

```shell
$ terraform import alicloud_cloud_storage_gateway_gateway_block_volume.example <gateway_id>:<index_id>
```