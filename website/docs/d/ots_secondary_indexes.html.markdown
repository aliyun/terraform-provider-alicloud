---
subcategory: "Table Store (OTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ots_secondary_index"
sidebar_current: "docs-alicloud-datasource-ots-secondary-index"
description: |-
  Provides a list of ots secondary index to the user.
---

# alicloud\_ots\_secondary\_index

This data source provides the ots secondary index of the current Alibaba Cloud user.

For information about OTS secondary index and how to use it, see [Secondary index overview](https://www.alibabacloud.com/help/en/tablestore/latest/secondary-index-overview).

-> **NOTE:** Available in v1.187.0+.

## Example Usage

``` terraform
data "alicloud_ots_secondary_indexes" "secondary_index_ds" {
  instance_name = "sample-instance"
  table_name = "sample-table"
  name_regex    = "sample-secondary-index"
  output_file   = "secondary-indexs.txt"
}

output "first_secondary_index_id" {
  value = "${data.alicloud_ots_secondary_indexes.secondary_index_ds.indexs.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required) The name of OTS instance.
* `table_name` - (Required) The name of OTS table.
* `ids` - (Optional) A list of secondary index IDs.
* `name_regex` - (Optional) A regex string to filter results by secondary index name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of secondary index IDs.
* `names` - A list of secondary index  names. 
* `indexes` - A list of indexes. Each element contains the following attributes:
    * `id` - The resource ID. The value is `<instance_name>:<table_name>:<indexName>:<indexType>`.
    * `instance_name` - The OTS instance name.
    * `table_name` - The table name of the OTS which could not be changed.
    * `index_name` - The index name of the OTS Table which could not be changed.
    * `index_type` - The index type of the OTS Table which could not be changed.
    * `primary_keys` - A list of primary keys for index, referenced from Table's primary keys or predefined columns.
    * `defined_columns` - A list of defined column for index, referenced from Table's primary keys or predefined columns.
   