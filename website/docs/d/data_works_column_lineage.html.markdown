---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_column_lineage"
sidebar_current: "docs-alicloud-datasource-data-works-column-lineage"
description: |-
  Provides the column lineage of a Data Works column to the user.
---

# alicloud_data_works_column_lineage

This data source provides the column lineage of a Data Works column of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.295.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_data_works_column_lineage" "default" {
  column_guid = "odps.project_name.table_name.column_name"
  direction   = "down"
}

output "first_lineage_column_guid" {
  value = data.alicloud_data_works_column_lineage.default.column_lineages.0.column_guid
}
```

Querying upstream lineage by data source type and table

```terraform
data "alicloud_data_works_column_lineage" "default" {
  data_source_type = "odps"
  database_name    = "your_database"
  table_name       = "your_table"
  column_name      = "your_column"
  direction        = "up"
}
```

## Argument Reference

The following arguments are supported:

* `column_guid` - (Optional) The globally unique identifier (GUID) of the column whose lineage is queried. It is one of the ways to locate the column.
* `direction` - (Required) The direction of the lineage to query. Valid values: `up` (upstream) and `down` (downstream).
* `data_source_type` - (Optional) The type of the data source. Valid values: `odps`, `emr`.
* `cluster_id` - (Optional) The ID of the EMR cluster. Required when the data source type is `emr`.
* `database_name` - (Optional) The name of the database.
* `table_name` - (Optional) The name of the table.
* `column_name` - (Optional) The name of the column.
* `ids` - (Optional, Computed) A list of Column GUIDs. The data source only returns the lineage entries whose column GUID matches one of the supplied values.
* `page_size` - (Optional) The number of items to return per page. The default value is `50` and the maximum value is `100`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

-> **NOTE:** The column to query can be located either by `column_guid` alone, or by the combination of `data_source_type`, `cluster_id`, `database_name`, `table_name` and `column_name`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `column_lineages` - A list of Data Works column lineage entries. Each element contains the following attributes:
  * `column_guid` - The globally unique identifier of the lineage column.
  * `column_name` - The name of the lineage column.
  * `table_name` - The name of the table that the lineage column belongs to.
  * `database_name` - The name of the database that the lineage column belongs to.
  * `cluster_id` - The ID of the EMR cluster that the lineage column belongs to.
  * `region_id` - The region ID of the data source.
