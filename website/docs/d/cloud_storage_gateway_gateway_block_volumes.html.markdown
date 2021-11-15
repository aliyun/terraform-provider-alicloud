---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_gateway_block_volumes"
sidebar_current: "docs-alicloud-datasource-cloud-storage-gateway-gateway-block-volumes"
description: |-
  Provides a list of Cloud Storage Gateway Gateway Block Volumes to the user.
---

# alicloud\_cloud\_storage\_gateway\_gateway\_block\_volumes

This data source provides the Cloud Storage Gateway Gateway Block Volumes of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.144.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cloud_storage_gateway_gateway_block_volumes" "ids" {
  gateway_id = "example_value"
  ids        = ["example_value-1", "example_value-2"]
}
output "cloud_storage_gateway_gateway_block_volume_id_1" {
  value = data.alicloud_cloud_storage_gateway_gateway_block_volumes.ids.volumes.0.id
}

data "alicloud_cloud_storage_gateway_gateway_block_volumes" "nameRegex" {
  gateway_id = "example_value"
  name_regex = "^my-GatewayBlockVolume"
}
output "cloud_storage_gateway_gateway_block_volume_id_2" {
  value = data.alicloud_cloud_storage_gateway_gateway_block_volumes.nameRegex.volumes.0.id
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, ForceNew) The Gateway ID.
* `ids` - (Optional, ForceNew, Computed)  A list of Gateway Block Volume IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Gateway Block Volume name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of volume. Valid values:
    - `0`: Normal condition.
    - `1`: Failed to create volume.
    - `2`: Failed to delete volume.
    - `3`: Failed to enable target.
    - `4`: Failed to disable target.
    - `5`: Database error.
    - `6`: Failed to enable cache.
    - `7`: Failed to disable cache.
    - `8`: System error.


## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Gateway Block Volume names.
* `volumes` - A list of Cloud Storage Gateway Gateway Block Volumes. Each element contains the following attributes:
	* `address` - The IP ADDRESS.
	* `cache_mode` - The Block volume set mode to cache mode. Value values: `Cache`, `WriteThrough`.
	* `chap_enabled` - Whether to enable iSCSI access of CHAP authentication, which currently supports both CHAP inbound authentication.  Default value: `false`.
	* `chap_in_user` - The Inbound CHAP user.**NOTE:** When the `chap_enabled` is  `true` is,The `chap_in_password` is valid.
	* `chunk_size` - The Block volume storage allocation unit.  Valid values: `8192`, `16384`, `32768`, `65536`, `131072`. Default value: `32768`. Unit: `Byte`.
	* `disk_id` - The cache disk ID.
	* `disk_type` - The cache disk type.
	* `enabled` - Whether to enable Volume.
	* `gateway_block_volume_name` - The Block volume name.  The name must be 1 to 32 characters in length, and can contain lowercase letters, numbers.
	* `gateway_id` - The Gateway ID.
	* `id` - The ID of the Gateway Block Volume. The value formats as `<gateway_id>:<index_id>`.
	* `index_id` - The ID of the index.
	* `local_path` - CThe Cache disk to local path. **NOTE:**  When the `cache_mode` is  `Cache` is,The `chap_in_password` is valid.
	* `lun_id` - The Lun identifier.
	* `operation_state` - The operation state.
	* `oss_bucket_name` - The name of the OSS Bucket.
	* `oss_bucket_ssl` - Whether to enable SSL access your OSS Buckets. Default value: `true`.
	* `oss_endpoint` - The endpoint of the OSS Bucket.
	* `port` - The Port.
	* `protocol` - The Protocol.
	* `size` - The Volume size.
	* `state` - The Buffer status.
	* `status` - The status of volume.
	* `target` - The target.
	* `total_download` - The total amount of downloaded data. Unit: `B`.
	* `total_upload` - The total amount of uploaded data. Unit: `B`.