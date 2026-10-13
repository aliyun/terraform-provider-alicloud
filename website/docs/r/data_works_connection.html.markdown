---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_connection"
description: |-
  Provides a Alicloud Data Works Connection resource.
---

# alicloud_data_works_connection

Provides a Data Works Connection resource.

Data Works Connection represents a data source connection registered under a Data Works project. It captures the connection type, environment, content (connection string), and lifecycle status used by DataWorks data integration tasks.

For information about Data Works Connection and how to use it, see [What is Connection](https://www.alibabacloud.com/help/en/dataworks/developer-reference/api-dataworks-public-2020-05-18-listconnections).

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
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, ForceNew, Int) The ID of the Data Works project to which the connection belongs.
* `connection_name` - (Required, ForceNew) The name of the connection. The name must be unique within the project environment.
* `connection_type` - (Required, ForceNew) The type of the connection, such as `mysql`, `odps`, `emr`, or other supported data source types.
* `sub_type` - (Optional, ForceNew) The sub type of the connection, used to further distinguish the data source category.
* `env_type` - (Required, Int) The environment type of the connection. Valid values: `0` (development environment), `1` (production environment).
* `content` - (Required, Sensitive) The connection content in JSON format. It contains the connection string parameters such as host, port, database, username, and password.
* `description` - (Optional) The description of the connection.
* `status` - (Optional, Computed) The status of the connection. Valid values: `1` (active), `2` (disabled).

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of the connection. The format is `<project_id>:<connection_id>`.
* `connection_id` - (Int) The ID of the connection.
* `operator` - The operator who created or modified the connection.
* `connect_status` - (Int) The connectivity status of the connection.
* `binding_calc_engine_id` - (Int) The ID of the bound compute engine.
* `gmt_modified` - The last modified time of the connection.
* `sequence` - (Int) The sequence of the connection.
* `shared` - (Bool) Whether the connection is shared.
* `default_engine` - (Bool) Whether the connection is the default engine.
* `create_time` - The creation time of the connection.
* `tenant_id` - (Int) The tenant ID of the connection.
* `region_id` - The region ID of the connection.

## Timeouts

The `timeouts` block allows you to specify [certain time limits](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) for individual operations:

* `create` - (Defaults to 5 mins) Used when creating the connection.
* `update` - (Defaults to 5 mins) Used when updating the connection.
* `delete` - (Defaults to 5 mins) Used when deleting the connection.

## Import

Data Works Connection can be imported using the id, e.g.

```shell
terraform import alicloud_data_works_connection.default <project_id>:<connection_id>
```
