---
subcategory: "DLF"
layout: "alicloud"
page_title: "Alicloud: alicloud_dlf_catalog"
sidebar_current: "docs-alicloud-resource-dlf-catalog"
description: |-
  Provides a Data Lake Formation (DlfNext) Catalog resource.
---

# alicloud_dlf_catalog

Provides a Data Lake Formation (DlfNext) Catalog resource. A catalog is the top-level container of a Data Lake Formation deployment; databases and tables live under it.

For information about the DlfNext API, see [Data Lake Formation API reference](https://www.alibabacloud.com/help/en/dlf/dlf-2-0/developer-reference/api-dlfnext-2025-03-10-dir/).

-> **NOTE:** Available since v2.3.0.

## Example Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

resource "alicloud_dlf_catalog" "default" {
  name      = "tf-dlf-catalog"
  type      = "PAIMON"
  is_shared = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the catalog. It is the path identifier used by Get/Alter/Drop and is unique within the account.
* `type` - (Required, ForceNew) The type of the catalog. Valid values: `PAIMON`, `ICEBERG`.
* `is_shared` - (Optional, ForceNew) Whether the catalog is shared. Defaults to `false`.
* `share_id` - (Optional, ForceNew) The share ID. Required when `is_shared` is `true`.

  -> **NOTE:** The catalog `options` are not settable. DlfNext populates storage defaults (such as `dlf.trashed-file-retained-days` and `storage.data.redundancy.type`) on its own; surfacing them as a user-set map would conflict with these server-managed keys. `options` is therefore exported read-only (see Attributes Reference).

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in Terraform, equal to the catalog `name`.
* `catalog_id` - The server-generated catalog ID (for example `clg-paimon-xxxx`). Use it as `catalog_id` of `alicloud_dlf_database` and `alicloud_dlf_table`.
* `options` - The extended options of the catalog, a map of key-value strings. Read-only; DlfNext populates storage defaults.
* `status` - The catalog status: `NEW`, `INITIALIZING`, `INITIALIZE_FAILED`, `RUNNING`, `TERMINATED`, `DELETING`, `DELETE_FAILED`, `DELETED`, `STORAGE_RESTRICTED`. Create waits for `RUNNING`.
* `owner` - The owner principal of the catalog.
* `created_at` - The creation time, epoch milliseconds.
* `created_by` - The creator principal of the catalog.
* `updated_at` - The last modification time, epoch milliseconds.
* `updated_by` - The last modifier principal of the catalog.

## Import

Catalog can be imported using the id, e.g.

```shell
terraform import alicloud_dlf_catalog.default <catalog_name>
```
