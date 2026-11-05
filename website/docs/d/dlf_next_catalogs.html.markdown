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

-> **NOTE:** Available since v2.3.0.

## Example Usage

Basic Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

resource "alicloud_dlf_next_catalog" "default" {
  name = "tf-demo-catalog"
  type = "PAIMON"
}

data "alicloud_dlf_next_catalogs" "default" {
  catalog_name_pattern = alicloud_dlf_next_catalog.default.name
}

output "first_catalog_name" {
  value = data.alicloud_dlf_next_catalogs.default.catalogs.0.name
}
```

## Argument Reference

The following arguments are supported:

* `catalog_name_pattern` - (Optional) A pattern to filter catalog names.
* `name_regex` - (Optional) A regex string to filter results by catalog name.
* `ids` - (Optional, Computed) A list of catalog names to filter results. The matched catalog names are also returned in `ids`.
* `output_file` - (Optional) File name where to save data source results after running `terraform plan`.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of matched catalog names.
* `names` - A list of matched catalog names.
* `catalogs` - A list of DLF Next Catalogs. Each element contains the following attributes:
  * `name` - The name of the catalog.
  * `type` - The type of the catalog.
  * `options` - The configuration options of the catalog.
  * `is_shared` - Whether the catalog is shared.
  * `share_id` - The share ID of the catalog.
  * `id` - The catalog name, which is also the Terraform resource id of `alicloud_dlf_next_catalog`.
  * `region_id` - The region ID of the catalog.
  * `status` - The status of the catalog. Valid values: `NEW`, `INITIALIZING`, `INITIALIZE_FAILED`, `RUNNING`, `TERMINATED`, `DELETING`, `DELETE_FAILED`, `DELETED`, `STORAGE_RESTRICTED`.
  * `owner` - The owner of the catalog.
  * `created_at` - The creation time of the catalog, in milliseconds since the epoch, as a string.
  * `created_by` - The creator of the catalog.
  * `updated_at` - The update time of the catalog, in milliseconds since the epoch, as a string.
  * `updated_by` - The updater of the catalog.

## RAM Permissions

The data source calls the DLF Next API, which is authorized at the API (RAM) layer. The provider does not automatically grant any RAM permissions to the execution identity (the AccessKey, STS token or RAM role used to run Terraform); the credentials used for `terraform plan`/`apply` must already hold the actions below.

* `dlf:ListCatalogs` - list catalogs (required to read the data source). The exact action set follows the current DLF RAM / CLI official documentation; `dlf:ListCatalogs` is the minimum required action.
