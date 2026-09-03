---
subcategory: "Realtime Compute for Apache Flink"
layout: "alicloud"
page_title: "Alicloud: alicloud_realtime_compute_sql_files"
sidebar_current: "docs-alicloud-datasource-realtime-compute-sql-files"
description: |-
  Provides a list of Realtime Compute for Apache Flink SQL Files to the user.
---

# alicloud_realtime_compute_sql_files

This data source provides a list of Realtime Compute for Apache Flink SQL Files in an Alibaba Cloud account.

For information about Realtime Compute for Apache Flink and how to use it, refer to [Realtime Compute for Apache Flink](https://www.alibabacloud.com/help/en/flink/developer-reference/api-ververica-2022-07-18).

-> **NOTE:** Available since v1.292.0.

## Example Usage

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_realtime_compute_sql_files" "default" {
  workspace = "your-workspace-id"
  namespace = "your-namespace"
  ids       = ["<workspace>:<namespace>:<sql_file_id>"]
}

output "sql_files" {
  value = data.alicloud_realtime_compute_sql_files.default.sql_files
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) The ID of the workspace.
* `namespace` - (Required) The name of the namespace.
* `ids` - (Optional) A list of SQL File IDs. The ID is formatted as `<workspace>:<namespace>:<sql_file_id>`.
* `name_regex` - (Optional) A regex string to filter results by the SQL file name.
* `enable_details` - (Optional) Default to false. Set true to show more attributes of the SQL files.
* `output_file` - (Optional) The name of output file that saves the filter results.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of SQL File IDs.
* `names` - A list of SQL file names.
* `sql_files` - A list of Realtime Compute SQL Files. Each element contains the following attributes:
  * `id` - The ID of the SQL file, formatted as `<workspace>:<namespace>:<sql_file_id>`.
  * `workspace` - The ID of the workspace.
  * `namespace` - The name of the namespace.
  * `sql_file_id` - The ID of the SQL file.
  * `name` - The name of the SQL file.
  * `sql_script` - The SQL script content of the SQL file.
  * `batch_mode` - The batch mode of the SQL file.
  * `session_cluster_name` - The name of the session cluster.
  * `parent_id` - The parent resource ID of the SQL file.
  * `description` - The description of the SQL file.
