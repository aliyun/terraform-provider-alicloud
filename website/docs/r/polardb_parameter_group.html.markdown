---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_parameter_group"
description: |-
  Provides a Alicloud Polar Db Parameter Group resource.
---

# alicloud_polardb_parameter_group

Provides a Polar Db Parameter Group resource.



For information about Polar Db Parameter Group and how to use it, see [What is Parameter Group](https://www.alibabacloud.com/help/en/polardb/polardb-for-mysql/user-guide/apply-a-parameter-template).

-> **NOTE:** Available since v1.183.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_polardb_parameter_group" "example" {
  parameter_group_name = "example_value"
  db_type              = "MySQL"
  db_version           = "8.0"
  parameters {
    param_name  = "wait_timeout"
    param_value = "86400"
  }
  description = "example_value"
}
```

## Argument Reference

The following arguments are supported:
* `db_type` - (Required, ForceNew) The type of the database engine. Only `MySQL` is supported.
* `db_version` - (Required, ForceNew) The version of the database engine. Valid values: 
  - **5.6** 
  - **5.7** 
  - **8.0**
* `description` - (Optional, ForceNew) ParameterGroupDesc
* `parameter_group_name` - (Optional, ForceNew, Available since v1.263.0) The name of the parameter template. The name must meet the following requirements:

  - It must start with a letter and can contain letters, digits, and underscores (_). It cannot contain Chinese characters or end with an underscore (_).

  - It must be 8 to 64 characters in length.
* `parameters` - (Required, ForceNew, Set) ParameterDetail See [`parameters`](#parameters) below.
* `name` - (Deprecated since v1.263.0). Field 'name' has been deprecated from provider version 1.263.0. New field 'parameter_group_name' instead.

### `parameters`

The parameters supports the following, you can view all parameter details for the target database engine version database cluster through the [DescribeParameterTemplates](https://next.api.alibabacloud.com/document/polardb/2017-08-01/DescribeParameterTemplates), including parameter name, value.
* `param_value` - (Optional, ForceNew) param value
* `param_name` - (Optional, ForceNew) param name

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Parameter Group.
* `delete` - (Defaults to 5 mins) Used when delete the Parameter Group.

## Import

Polar Db Parameter Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_parameter_group.example <id>
```