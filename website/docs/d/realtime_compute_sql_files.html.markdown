---
subcategory: "Realtime Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_realtime_compute_sql_files"
sidebar_current: "docs-alicloud-datasource-realtime-compute-sql-files"
description: |-
  Provides a list of Realtime Compute Sql File owned by an Alibaba Cloud account.
---

# alicloud_realtime_compute_sql_files

This data source provides Realtime Compute Sql File available to the user.[What is Sql File](https://next.api.alibabacloud.com/document/ververica/2022-07-18/CreateSqlFile)

-> **NOTE:** Available since v1.292.0.

## Example Usage

```terraform
variable "name" {
  default = "tfexample-sql-file"
}

data "alicloud_realtime_compute_sql_files" "default" {
  workspace = "a14bda1c4a****"
  namespace = "default-namespace"
  ids       = ["a14bda1c4a****:default-namespace:123456"]
}

output "first_sql_file_id" {
  value = data.alicloud_realtime_compute_sql_files.default.files.0.id
}
```

## Argument Reference

The following arguments are supported:
* `namespace` - (Required) The name of the namespace.
* `workspace` - (Required) The ID of the workspace.
* `ids` - (Optional, Computed) A list of Sql File IDs. The value is formulated as `<workspace>:<namespace>:<sql_file_id>`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Sql File IDs.
* `files` - A list of Sql File Entries. Each element contains the following attributes:
  * `batch_mode` - Whether the SQL query script runs in batch mode.
  * `description` - The description of the SQL file.
  * `name` - The name of the SQL file.
  * `namespace` - The name of the namespace.
  * `parent_id` - The ID of the parent folder of the SQL file.
  * `session_cluster_name` - The name of the session cluster that runs the SQL query script.
  * `sql_file_id` - The ID of the SQL file.
  * `sql_script` - The SQL script content.
  * `workspace` - The ID of the workspace.
  * `id` - The ID of the SQL file. The value is formulated as `<workspace>:<namespace>:<sql_file_id>`.
