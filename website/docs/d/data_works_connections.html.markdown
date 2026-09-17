---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_connections"
sidebar_current: "docs-alicloud-datasource-data-works-connections"
description: |-
  Provides a list of Data Works Connections to the user.
---

# alicloud_data_works_connections

This data source provides the Data Works Connections of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.238.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}

resource "alicloud_data_works_connection" "default" {
  project_id      = 12345
  connection_name = var.name
  connection_type = "mysql"
  sub_type        = "mysql"
  env_type        = 0
  content = jsonencode({
    database = "tf_example_db"
    host     = "127.0.0.1"
    password = "tf_example_pw"
    port     = "3306"
    username = "tf_example_user"
  })
  description = "tf-example-connection-desc"
}

data "alicloud_data_works_connections" "default" {
  project_id      = alicloud_data_works_connection.default.project_id
  connection_type = "mysql"
  name            = var.name
  ids             = [alicloud_data_works_connection.default.connection_id]
}

output "connection_id" {
  value = data.alicloud_data_works_connections.default.connections.0.connection_id
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, Int) The ID of the Data Works project to which the connections belong.
* `connection_type` - (Optional) The type of the connection. Used to filter connections by type.
* `env_type` - (Optional, Int) The environment type of the connection. Valid values: `0` (development), `1` (production).
* `ids` - (Optional, Computed) A list of Connection IDs. Used to filter connections by specific IDs.
* `name` - (Optional) The name of the connection. Used to filter connections by name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional) The status of the connection. Valid values: `1` (active), `2` (disabled).
* `sub_type` - (Optional) The sub type of the connection. Used to filter connections by sub type.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `connections` - A list of Data Works Connections. Each element contains the following attributes:
  * `id` - The ID of the connection. The format is `<project_id>:<connection_id>`.
  * `binding_calc_engine_id` - (Int) The ID of the bound compute engine.
  * `connection_id` - (Int) The ID of the connection.
  * `connection_name` - The name of the connection.
  * `connection_type` - The type of the connection.
  * `connect_status` - (Int) The connectivity status of the connection.
  * `content` - (Sensitive) The connection content in JSON format.
  * `create_time` - The creation time of the connection.
  * `default_engine` - (Bool) Whether the connection is the default engine.
  * `description` - The description of the connection.
  * `env_type` - (Int) The environment type of the connection.
  * `gmt_modified` - The last modified time of the connection.
  * `operator` - The operator who created or modified the connection.
  * `project_id` - (Int) The ID of the project.
  * `region_id` - The region ID of the connection.
  * `sequence` - (Int) The sequence of the connection.
  * `shared` - (Bool) Whether the connection is shared.
  * `status` - The status of the connection.
  * `sub_type` - The sub type of the connection.
  * `tenant_id` - (Int) The tenant ID of the connection.
