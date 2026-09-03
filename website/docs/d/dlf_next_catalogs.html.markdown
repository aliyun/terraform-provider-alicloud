---
subcategory: "dlf_next"
layout: "alicloud"
page_title: "Alicloud: alicloud_dlf_next_catalogs"
sidebar_current: "docs-alicloud-datasource-dlf-next-catalogs"
description: |-
  Provides a list of DLF Next Catalogs to the user.
---

# alicloud_dlf_next_catalogs

This data source provides the DLF Next Catalogs of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.232.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_dlf_next_catalogs" "default" {
  catalog_name_pattern = "tf-test-*"
}

output "first_catalog_name" {
  value = data.alicloud_dlf_next_catalogs.default.catalogs.0.name
}
```

## Argument Reference

The following arguments are supported:

* `catalog_name_pattern` - (Optional) A pattern to filter catalog names.
* `name_regex` - (Optional) A regex string to filter results by catalog name.
* `ids` - (Optional) A list of catalog names to filter results.
* `output_file` - (Optional) File name where to save data source results after running `terraform plan`.

## Attributes Reference

The following attributes are exported:

* `catalogs` - A list of DLF Next Catalogs. Each element contains the following attributes:
  * `name` - The name of the catalog.
  * `type` - The type of the catalog.
  * `options` - The configuration options of the catalog.
  * `is_shared` - Whether the catalog is shared.
  * `share_id` - The share ID of the catalog.
  * `id` - The ID of the catalog.
  * `region_id` - The region ID of the catalog.
  * `status` - The status of the catalog.
  * `owner` - The owner of the catalog.
  * `created_at` - The creation time of the catalog.
  * `created_by` - The creator of the catalog.
  * `updated_at` - The update time of the catalog.
  * `updated_by` - The updater of the catalog.
