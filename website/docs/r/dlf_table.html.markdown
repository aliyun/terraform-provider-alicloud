---
subcategory: "DLF"
layout: "alicloud"
page_title: "Alicloud: alicloud_dlf_table"
sidebar_current: "docs-alicloud-resource-dlf-table"
description: |-
  Provides a Data Lake Formation (DlfNext) Table resource.
---

# alicloud_dlf_table

Provides a Data Lake Formation (DlfNext) Table resource. A table lives under a database and defines a schema of typed columns. Changing any argument of a table forces replacement (destroy and recreate), as DlfNext does not expose an in-place alter for table structure.

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
}

resource "alicloud_dlf_table" "default" {
  catalog_id    = alicloud_dlf_catalog.default.catalog_id
  database_name = alicloud_dlf_database.default.name
  name          = "tf-dlf-table"
  comment       = "managed by terraform"

  # `fields` is a repeated block. Each block has `name`, a SQL-like `type`
  # and an optional `description`. Changing any field forces replacement.
  # Partition key columns (see `partition_keys` below) must also be declared
  # as a `fields` block.
  fields {
    name = "id"
    type = "BIGINT"
  }

  fields {
    name = "name"
    type = "STRING"
  }

  fields {
    name        = "tags"
    type        = "MAP<STRING,STRING>"
    description = "tag map"
  }

  fields {
    name = "dt"
    type = "STRING"
  }

  partition_keys = ["dt"]
  primary_keys   = ["id"]
}
```

## Argument Reference

The following arguments are supported:

* `catalog_id` - (Required, ForceNew) The ID of the parent catalog (the `catalog_id` of an `alicloud_dlf_catalog` resource).
* `database_name` - (Required, ForceNew) The name of the parent database (the `name` of an `alicloud_dlf_database` resource).
* `name` - (Required, ForceNew) The name of the table.
* `comment` - (Optional, ForceNew) A free-form comment for the table.
* `fields` - (Optional, ForceNew) The table columns. See [`fields`](#fields) below.
* `options` - (Optional, ForceNew) Extended table options, a map of key-value strings. Also exported (Computed) so the server-reported options round-trip into state.
* `partition_keys` - (Optional, ForceNew) The partition key column names. Each partition key column must also be declared as a `fields` block.
* `primary_keys` - (Optional, ForceNew) The primary key column names.

### fields

The `fields` block supports:

* `name` - (Required, ForceNew) The column name.
* `type` - (Required, ForceNew) The column type, a SQL-like expression. Primitive types such as `STRING`, `BIGINT`, `DOUBLE`, `BOOLEAN` are written as-is. Parametric types use angle brackets: `ARRAY<INT>`, `MULTISET<STRING>`, `MAP<STRING,BIGINT>`, and `ROW<id INT, name STRING>`.
* `description` - (Optional, ForceNew) An optional description of the column.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in Terraform, in the form `<catalog_id>:<database_name>:<name>`.
* `table_id` - The server-generated table ID.
* `options` - The extended table options, a map of key-value strings (Computed).
* `owner` - The owner principal of the table.
* `status` - The table status.
* `created_at` - The creation time, epoch milliseconds.
* `created_by` - The creator principal of the table.
* `updated_at` - The last modification time, epoch milliseconds.
* `updated_by` - The last modifier principal of the table.

## Import

Table can be imported using the id, e.g.

```shell
terraform import alicloud_dlf_table.default <catalog_id>:<database_name>:<table_name>
```
