---
subcategory: "Table Store (OTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ots_search_index"
sidebar_current: "docs-alicloud-datasource-ots-search-index"
description: |-
  Provides a list of ots search index to the user.
---

# alicloud\_ots\_search\_index

This data source provides the ots search index of the current Alibaba Cloud user.

For information about OTS search index and how to use it, see [Search index overview](https://www.alibabacloud.com/help/en/tablestore/latest/search-index-overview).

-> **NOTE:** Available in v1.187.0+.

## Example Usage

```
data "alicloud_ots_search_indexes" "search_index_ds" {
  instance_name = "sample-instance"
  table_name = "sample-table"
  name_regex    = "sample-search-index"
  output_file   = "search-indexs.txt"
}

output "first_search_index_id" {
  value = "${data.alicloud_ots_search_indexes.search_index_ds.indexs.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required) The name of OTS instance.
* `table_name` - (Required) The name of OTS table.
* `ids` - (Optional) A list of search index IDs.
* `name_regex` - (Optional) A regex string to filter results by search index name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of search index IDs.
* `names` - A list of search index  names. 
* `indexes` - A list of indexes. Each element contains the following attributes:
    * `id` - The resource ID. The value is `<instance_name>:<table_name>:<indexName>:<indexType>`.
    * `instance_name` - The OTS instance name.
    * `table_name` - The table name of the OTS which could not be changed.
    * `index_name` - The index name of the OTS Table which could not be changed.
    * `create_time` - The creation time of the index.
    * `time_to_live` - TTL of index.
    * `sync_phase` - The synchronization state of the index.
    * `current_sync_timestamp` - Timestamp for sync phase.
    * `storage_size` - Storage space occupied by index.
    * `row_count` - The number of rows of data for index.
    * `reserved_read_cu` - Reserve related resources for the index.
    * `metering_last_update_time` - Last update time for metering data..
    * `schema` - JSON representation of the schema of index.
   