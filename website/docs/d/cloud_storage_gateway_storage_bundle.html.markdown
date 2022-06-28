---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_storage_bundles"
sidebar_current: "docs-alicloud-datasource-cloud-storage-gateway-storage-bundles"
description: |-
  Provides a list of Cloud Storage Gateway Storage Bundles to the user.
---

# alicloud\_cloud\_storage\_gateway\_storage\_bundles

This data source provides the Cloud Storage Gateway Storage Bundles of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.116.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cloud_storage_gateway_storage_bundles" "example" {
  backend_bucket_region_id = "cn-beijing"
  ids                      = ["sb-0008xeww8yg2********"]
  name_regex               = "the_resource_name"
}

output "first_cloud_storage_gateway_storage_bundle_id" {
  value = data.alicloud_cloud_storage_gateway_storage_bundles.example.bundles.0.id
}
```

## Argument Reference

The following arguments are supported:

* `backend_bucket_region_id` - (Required, ForceNew) The backend bucket region id.
* `ids` - (Optional, ForceNew, Computed)  A list of Storage Bundle IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Storage Bundle name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Storage Bundle names.
* `bundles` - A list of Cloud Storage Gateway Storage Bundles. Each element contains the following attributes:
	* `description` - storage bundle description.
	* `id` - The ID of the Storage Bundle.
	* `location` - The location of storage bundle.
	* `storage_bundle_id` - The id of storage bundle.
	* `storage_bundle_name` - The name of storage bundle.
	* `create_time` - The created time of storage bundle.