---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_storage_bundles"
sidebar_current: "docs-alicloud-datasource-cloud-storage-gateway-storage-bundles"
description: |-
  Provides a list of Cloud Storage Gateway Storage Bundles to the user.
---

# alicloud_cloud_storage_gateway_storage_bundles

This data source provides the Cloud Storage Gateway Storage Bundles of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.132.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_cloud_storage_gateway_storage_bundles" "example" {
  backend_bucket_region_id = "cn-hangzhou"
  name_regex               = "^my-StorageBundle"
}

output "storage_bundle_id" {
  value = data.alicloud_cloud_storage_gateway_storage_bundles.example.bundles.0.storage_bundle_id
}
```

## Argument Reference

The following arguments are supported:

* `backend_bucket_region_id` - (Required, ForceNew) The region id of the backend bucket.
* `ids` - (Optional, ForceNew, Computed) A list of Storage Bundle IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Storage Bundle name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `page_number` - (Optional, ForceNew) Page number.
* `page_size` - (Optional, ForceNew) Page size. Default to 50.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Storage Bundle IDs.
* `names` - A list of Storage Bundle names.
* `total_count` - Total count of Storage Bundle.
* `bundles` - A list of Cloud Storage Gateway Storage Bundles. Each element contains the following attributes:
  * `description` - Storage Bundle description.
  * `location` - Storage Bundle location.
  * `id` - The ID of the Storage Bundle.
  * `storage_bundle_id` - The ID of the Storage Bundle.
  * `storage_bundle_name` - The name of the Storage Bundle.
  * `create_time` - Storage Bundle created time.
