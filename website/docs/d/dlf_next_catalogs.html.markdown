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

-> **NOTE:** Available since v1.293.0.

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
* `names` - (Optional) A list of catalog names to filter results.
* `output_file` - (Optional) File name where to save data source results after running `terraform plan`.

## Attributes Reference

The following attributes are exported:

* `catalogs` - A list of DLF Next Catalogs. Each element contains the following attributes:
  * `name` - The name of the catalog.
  * `type` - The type of the catalog.
  * `options` - The configuration options of the catalog.
  * `is_shared` - Whether the catalog is shared.
  * `share_id` - The share ID of the catalog.
  * `id` - The server-generated id of the catalog (of the form `clg-xxx`). This is distinct from the catalog name, which is the Terraform resource id of `alicloud_dlf_next_catalog`.
  * `region_id` - The region ID of the catalog.
  * `status` - The status of the catalog. Valid values: `NEW`, `INITIALIZING`, `INITIALIZE_FAILED`, `RUNNING`, `TERMINATED`, `DELETING`, `DELETE_FAILED`, `DELETED`, `STORAGE_RESTRICTED`.
  * `owner` - The owner of the catalog.
  * `created_at` - The creation time of the catalog.
  * `created_by` - The creator of the catalog.
  * `updated_at` - The update time of the catalog.
  * `updated_by` - The updater of the catalog.

## RAM Permissions

The data source calls the DLF Next API, which is authorized at the API (RAM) layer. The provider does not automatically grant any RAM permissions to the execution identity (the AccessKey, STS token or RAM role used to run Terraform); the credentials used for `terraform plan`/`apply` must already hold the actions below.

* `dlf:ListCatalogs` - list catalogs (required to read the data source). The exact action set follows the current DLF RAM / CLI official documentation; `dlf:ListCatalogs` is the minimum required action.

-> **NOTE:** Catalog CRUD / List acceptance tests only verify the API-layer (RAM) authorization for DLF Next Catalog management. They do not verify DLF data-layer authorization to access table data through Paimon REST, Iceberg REST or HMS, so data-access permissions cannot be claimed as verified by running this data source.
