---
subcategory: "DLF"
layout: "alicloud"
page_title: "Alicloud: alicloud_dlf_database"
sidebar_current: "docs-alicloud-resource-dlf-database"
description: |-
  Provides a Data Lake Formation (DlfNext) Database resource.
---

# alicloud_dlf_database

Provides a Data Lake Formation (DlfNext) Database resource. A database lives under a catalog and contains tables.

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

resource "alicloud_dlf_database" "default" {
  catalog_id = alicloud_dlf_catalog.default.catalog_id
  name       = "tf-dlf-database"
  options = {
    key1 = "value1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `catalog_id` - (Required, ForceNew) The ID of the parent catalog (the `catalog_id` of an `alicloud_dlf_catalog` resource).
* `name` - (Required, ForceNew) The name of the database.
* `options` - (Optional) The extended options of the database, a map of key-value strings. Omitting the block clears user-set keys. Updatable in place via AlterDatabase. The attribute is also exported (Computed) so the server-reported options round-trip into state.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in Terraform, in the form `<catalog_id>:<name>`.
* `database_id` - The server-generated database UUID.
* `location` - The storage location of the database.
* `options` - The extended options of the database, a map of key-value strings (Computed).
* `owner` - The owner principal of the database.
* `table_count` - The number of tables in the database.
* `status` - The database status.
* `created_at` - The creation time, epoch milliseconds.
* `created_by` - The creator principal of the database.
* `updated_at` - The last modification time, epoch milliseconds.
* `updated_by` - The last modifier principal of the database.

## Import

Database can be imported using the id, e.g.

```shell
terraform import alicloud_dlf_database.default <catalog_id>:<database_name>
```
