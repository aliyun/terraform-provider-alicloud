---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_data_source"
description: |-
  Provides a Alicloud Data Works Data Source resource.
---

# alicloud_data_works_data_source

Provides a Data Works Data Source resource.



For information about Data Works Data Source and how to use it, see [What is Data Source](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.237.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}

provider "alicloud" {
  region = "cn-chengdu"
}

resource "random_integer" "randint" {
  max = 999
  min = 1
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_data_works_project" "defaultkguw4R" {
  status                  = "Available"
  description             = "tf_desc"
  project_name            = "${var.name}${random_integer.randint.id}"
  pai_task_enabled        = "false"
  display_name            = "tf_new_api_display"
  dev_role_disabled       = "true"
  dev_environment_enabled = "false"
  resource_group_id       = data.alicloud_resource_manager_resource_groups.default.ids.0
}

resource "alicloud_data_works_data_source" "default" {
  type                       = "hive"
  data_source_name           = var.name
  connection_properties      = jsonencode({ "address" : [{ "host" : "127.0.0.1", "port" : "1234" }], "database" : "hive_database", "metaType" : "HiveMetastore", "metastoreUris" : "thrift://123:123", "version" : "2.3.9", "loginMode" : "Anonymous", "securityProtocol" : "authTypeNone", "envType" : "Prod", "properties" : { "key1" : "value1" } })
  connection_properties_mode = "UrlMode"
  project_id                 = alicloud_data_works_project.defaultkguw4R.id
  description                = var.name
}
```

## Argument Reference

The following arguments are supported:
* `connection_properties` - (Required, JsonString) Data source connection configuration information, including the connection address, access identity, and environment information. The data source environment EnvType information is a member property of this object, including DEV (Development Environment) and PROD (production environment). The value of EnvType is not case-sensitive.
* `connection_properties_mode` - (Required) The configuration mode of the data source. Different types of data sources have different configuration modes. For example, MySQL data sources support UrlMode and InstanceMode.
* `data_source_name` - (Required, ForceNew) The data source name. The name of a data source in a specific environment (development environment or production environment) is unique in a project.
* `description` - (Optional) Description of the data source
* `project_id` - (Required, ForceNew, Int) The ID of the project to which the data source belongs.
* `type` - (Required, ForceNew) The type of data source. For a list of data source types, see the values listed in the API documentation.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<project_id>:<data_source_id>`.
* `create_time` - The creation time of the resource
* `create_user` - Creator of the data source
* `data_source_id` - The first ID of the resource
* `modify_time` - Modification time
* `modify_user` - Modifier of the data source
* `qualified_name` - Business Unique Key of Data Source

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Data Source.
* `delete` - (Defaults to 5 mins) Used when delete the Data Source.
* `update` - (Defaults to 5 mins) Used when update the Data Source.

## Import

Data Works Data Source can be imported using the id, e.g.

```shell
$ terraform import alicloud_data_works_data_source.example <project_id>:<data_source_id>
```