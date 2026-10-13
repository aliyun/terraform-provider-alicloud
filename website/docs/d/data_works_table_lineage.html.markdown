---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_table_lineage"
sidebar_current: "docs-alicloud-datasource-data-works-table-lineage"
description: |-
  Provides a list of Data Works Table Lineage entries to the user.
---

# alicloud\_data\_works\_table\_lineage

This data source provides the Data Works Table Lineage of the current Alibaba Cloud user by direction (upstream or downstream) and optional table filters.

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
data "alicloud_data_works_table_lineage" "default" {
  direction        = "down"
  table_name       = "my_table"
  database_name    = "my_database"
  data_source_type = "odps"
}

output "first_lineage_table_guid" {
  value = data.alicloud_data_works_table_lineage.default.lineages.0.table_guid
}
```

## Argument Reference

The following arguments are supported:

* `direction` - (Required) The direction of the lineage. Valid values: `up` (upstream) and `down` (downstream).
* `table_guid` - (Optional) The unique identifier of the table. If specified, the lineage of the table identified by this GUID is returned.
* `table_name` - (Optional) The name of the table. Used together with `database_name` and `data_source_type` to identify a table when `table_guid` is unknown.
* `database_name` - (Optional) The name of the database where the table resides.
* `data_source_type` - (Optional) The data source type. Valid values: `odps` and `emr`. Required for EMR scenarios together with `cluster_id`.
* `cluster_id` - (Optional) The ID of the EMR cluster. This parameter is required for EMR scenarios.
* `page_size` - (Optional) The number of entries per page. Default value: `100`. Maximum value: `100`. The data source auto-paginates all results using `NextPrimaryKey`.
* `ids` - (Optional) A list of lineage ids used to filter the returned lineages. Only lineages whose `id` matches one of the given ids are returned.
* `output_file` - (Optional) File name where to write the result after running `terraform plan`.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of matching lineage ids. Each id is the `table_guid` of the lineage entry (or a composite `table_name:database_name:direction` fallback when the API does not return a `TableGuid`).
* `lineages` - A list of table lineage entries. Each element contains the following attributes:
  * `id` - The id of the lineage entry, same as `table_guid` when present.
  * `table_guid` - The unique identifier of the table.
  * `table_name` - The name of the table.
  * `database_name` - The name of the database.
  * `create_timestamp` - The creation timestamp of the table.
