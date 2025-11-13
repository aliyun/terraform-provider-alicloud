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
* `description` - (Optional, ForceNew) The description of the parameter template.
* `parameter_group_name` - (Optional, ForceNew, Available since v1.263.0) The name of the parameter template. The name must meet the following requirements:

  - It must start with a letter and can contain letters, digits, and underscores (_). It cannot contain Chinese characters or end with an underscore (_).

  - It must be 8 to 64 characters in length.
* `parameters` - (Required, ForceNew, Set) Details about the parameters. See [`parameters`](#parameters) below.

-> **NOTE:**  You can view all parameter details for the target database engine version database cluster through the [DescribeParameterTemplates](https://next.api.alibabacloud.com/document/polardb/2017-08-01/DescribeParameterTemplates), including parameter name, value.


The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated since v1.263.0). Field 'name' has been deprecated from provider version 1.263.0. New field 'parameter_group_name' instead.

### `parameters`

The parameters supports the following, you can view all parameter details for the target database engine version database cluster through the [DescribeParameterTemplates](https://next.api.alibabacloud.com/document/polardb/2017-08-01/DescribeParameterTemplates), including parameter name, value.
* `param_value` - (Optional, ForceNew) The value of the parameter.
* `param_name` - (Optional, ForceNew) The name of the parameter.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the parameter template was created. The time is in the yyyy-MM-ddTHH:mm:ssZ format. The time is displayed in UTC.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Parameter Group.
* `delete` - (Defaults to 5 mins) Used when delete the Parameter Group.

## Import

Polar Db Parameter Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_parameter_group.example <id>
```