---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_columns"
sidebar_current: "docs-alicloud-datasource-data-works-columns"
description: |-
  Provides a list of Data Works table columns to the user.
---

# alicloud_data_works_columns

This data source provides the Data Works table columns of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.294.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_data_works_columns" "default" {
  table_guid = "odps.project_name.table_name"
}

output "first_column_guid" {
  value = data.alicloud_data_works_columns.default.columns.0.column_guid
}
```

Querying by data source type and table name

```terraform
data "alicloud_data_works_columns" "default" {
  data_source_type = "maxcompute"
  database_name    = "your_database"
  table_name       = "your_table"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, Computed) A list of Column GUIDs. The data source only returns the columns whose GUID matches one of the supplied values.
* `table_guid` - (Optional) The globally unique identifier (GUID) of the table. It is one of the ways to locate the table.
* `data_source_type` - (Optional) The type of the data source. Valid values: `maxcompute`, `dlf`, `hms`, `holo`, `mysql`.
* `cluster_id` - (Optional) The ID of the cluster. Required when the data source type is `holo` or `mysql`.
* `database_name` - (Optional) The name of the database.
* `table_name` - (Optional) The name of the table.
* `page_size` - (Optional) The number of items to return per page. The default value is `50` and the maximum value is `100`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

-> **NOTE:** The table to query can be located either by `table_guid` alone, or by the combination of `data_source_type`, `cluster_id`, `database_name` and `table_name`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `columns` - A list of Data Works table columns. Each element contains the following attributes:
  * `column_guid` - The globally unique identifier of the column.
  * `column_name` - The name of the column.
  * `comment` - The comment of the column.
  * `column_type` - The type of the column.
  * `is_primary_key` - Whether the column is a primary key.
  * `region_id` - The region ID of the data source.
