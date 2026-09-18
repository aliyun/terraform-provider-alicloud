---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_node_io"
sidebar_current: "docs-alicloud-datasource-data-works-node-io"
description: |-
  Provides the Data Works Node Input/Output list of the current Alibaba Cloud user.
---

# alicloud_data_works_node_io

This data source provides the Data Works Node Input/Output list of the current Alibaba Cloud user according to the specified node, project environment and IO type.

-> **NOTE:** Available since v1.226.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_data_works_node_io" "default" {
  node_id     = "12345"
  project_env = "PROD"
  io_type     = "input"
}

output "node_io_table_names" {
  value = data.alicloud_data_works_node_io.default.node_ios[*].table_name
}
```

## Argument Reference

The following arguments are supported:

* `node_id` - (Required) The ID of the DataWorks node. ListNodeIO queries the upstream or downstream nodes of this node.
* `project_env` - (Required) The environment of the DataWorks project. Valid values: `PROD`, `DEV`.
* `io_type` - (Required) The type of the node input/output. It queries the upstream (`input`) or downstream (`output`) nodes only one level.
* `ids` - (Optional) A list of node IO entry IDs (the related table names). If set, only the entries whose table name matches one of the provided IDs are returned.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `node_ios` - A list of Data Works Node Input/Output entries. Each element contains the following attributes:
  * `id` - The ID of the node IO entry. It is the same value as `table_name`.
  * `table_name` - The name of the upstream or downstream table associated with the node.
  * `data` - The data content of the node input/output.
  * `node_id` - The ID of the DataWorks node.
