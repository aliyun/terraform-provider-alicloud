---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_dataworks_table"
description: |-
  Provides a Alicloud Data Works Table resource.
---

# alicloud_dataworks_table

Provides a Data Works Table resource.

For information about Data Works Table and how to use it, see [CreateTable](https://www.alibabacloud.com/help/en/dataworks/developer-reference/api-dataworks-public-2020-05-18-createtable).

-> **NOTE:** The table is created in the MaxCompute engine bound to the workspace. `table_guid` is the composite guid `odps.<engine>.<table_name>`; the plain `<table_name>` segment is used as the table name, and the full guid is used to read the table.

-> **NOTE:** Available since v1.295.0.

## Example Usage

Basic Usage

```terraform
variable "project_id" {
  default = "12345678" # a DataWorks workspace id that has a MaxCompute engine bound
}

variable "engine" {
  default = "my_odps_project" # the MaxCompute project (engine) bound to the workspace
}

resource "alicloud_dataworks_table" "default" {
  table_guid = "odps.${var.engine}.tf_example_table"
  project_id = var.project_id
  comment    = "tf example table"
  life_cycle = 30
  is_view    = 0
  intro_wiki = "intro content"
}
```

## Argument Reference

The following arguments are supported:

* `table_guid` - (Required, ForceNew) The composite guid of the table, in the form `odps.<engine>.<table_name>`. The plain `<table_name>` segment is used as the TableName of the create/update/delete API; the full guid is used to read the table.
* `project_id` - (Required, ForceNew) The ID of the DataWorks workspace that owns the table.
* `comment` - (Optional) The comment of the table.
* `life_cycle` - (Optional) The lifecycle of the table, in days.
* `location` - (Optional) The storage location of the external table.
* `owner_id` - (Optional) The owner ID of the table.
* `is_view` - (Optional) Whether the table is a view. `0`: table, `1`: view. Write-only: not returned by the read API.
* `category_id` - (Optional) The category ID of the table. Write-only: not returned by the read API.
* `physics_level_id` - (Optional) The physics level ID of the table. Write-only: not returned by the read API.
* `intro_wiki` - (Optional) The intro wiki content of the table. Drives the `UpdateMetaTableIntroWiki` API; write-only.
* `first_level_theme_id` - (Optional) The first-level theme ID. Drives the `UpdateTableModelInfo` API; write-only.
* `second_level_theme_id` - (Optional) The second-level theme ID. Drives the `UpdateTableModelInfo` API; write-only.
* `level_id` - (Optional) The model level ID. Drives the `UpdateTableModelInfo` API; write-only.
* `level_type` - (Optional) The model level type. Drives the `UpdateTableModelInfo` API; write-only.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, same as `table_guid`.
* `table_name` - The plain name of the table.
* `create_time` - The time when the table was created.
* `last_modify_time` - The time when the table was last modified.
* `tenant_id` - The tenant ID.
* `data_size` - The data size of the table.
* `is_visible` - Whether the table is visible.
* `region_id` - The region ID of the resource.

## Import

Data Works Table can be imported using the `id`, i.e. the composite `table_guid`, e.g.

```shell
$ terraform import alicloud_dataworks_table.default odps.my_odps_project.my_table
```
