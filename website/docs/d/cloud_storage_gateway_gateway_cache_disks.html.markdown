---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_gateway_cache_disks"
sidebar_current: "docs-alicloud-datasource-cloud-storage-gateway-gateway-cache-disks"
description: |-
  Provides a list of Cloud Storage Gateway Gateway Cache Disks to the user.
---

# alicloud\_cloud\_storage\_gateway\_gateway\_cache\_disks

This data source provides the Cloud Storage Gateway Gateway Cache Disks of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.144.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cloud_storage_gateway_gateway_cache_disks" "ids" {
  gateway_id = "example_value"
  ids        = ["example_value-1", "example_value-2"]
}
output "cloud_storage_gateway_gateway_cache_disk_id_1" {
  value = data.alicloud_cloud_storage_gateway_gateway_cache_disks.ids.disks.0.id
}

data "alicloud_cloud_storage_gateway_gateway_cache_disks" "status" {
  gateway_id = "example_value"
  ids        = ["example_value-1", "example_value-2"]
  status     = "0"
}
output "cloud_storage_gateway_gateway_cache_disk_id_2" {
  value = data.alicloud_cloud_storage_gateway_gateway_cache_disks.status.disks.0.id
}

```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, ForceNew) The ID of the gateway.
* `ids` - (Optional, ForceNew, Computed)  A list of Gateway Cache Disk IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `0`, `1`, `2`. `0`: Normal. `1`: Is about to expire. `2`: Has expired.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `disks` - A list of Cloud Storage Gateway Gateway Cache Disks. Each element contains the following attributes:
	* `cache_disk_category` - The category of eht cache disk.
	* `cache_disk_size_in_gb` - The size of the cache disk.
	* `cache_id` - The ID of the cache disk.
	* `expired_time` - The expiration time. Time stamp in seconds (s).
	* `gateway_id` - The ID of the gateway.
	* `id` - The ID of the Gateway Cache Disk.
	* `iops` - Per second of the input output.
	* `is_used` - Whether it is used.
	* `local_file_path` - The cache disk inside the device name.
	* `renew_url` - A renewal link of the cache disk.
	* `status` - The status of the resource.