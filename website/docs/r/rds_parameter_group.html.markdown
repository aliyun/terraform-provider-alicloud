---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_parameter_group"
sidebar_current: "docs-alicloud-resource-rds-parameter-group"
description: |-
  Provides a Alicloud RDS Parameter Group resource.
---

# alicloud\_rds\_parameter\_group

Provides a RDS Parameter Group resource.

For information about RDS Parameter Group and how to use it, see [What is Parameter Group](https://www.alibabacloud.com/help/en/doc-detail/144839.htm).

-> **NOTE:** Available in v1.119.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_rds_parameter_group" "default" {
  engine         = "mysql"
  engine_version = "5.7"
  param_detail {
    param_name  = "back_log"
    param_value = "4000"
  }
  param_detail {
    param_name  = "wait_timeout"
    param_value = "86460"
  }
  parameter_group_desc = "test"
  parameter_group_name = "test1234"
}
```

## Argument Reference

The following arguments are supported:

* `engine` - (Required, ForceNew) The database engine. Valid values: `mysql`, `mariadb`.
* `engine_version` - (Required, ForceNew) The version of the database engine. Valid values: mysql: `5.1`, `5.5`, `5.6`, `5.7`, `8.0`; mariadb: `10.3`.
* `param_detail` - (Required) Parameter list.
* `parameter_group_desc` - (Optional) The description of the parameter template.
* `parameter_group_name` - (Required) The name of the parameter template.

#### Block parameter_detail

The param_detail supports the following: 

* `param_name` - (Required) The name of a parameter.
* `param_value` - (Required) The value of a parameter.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Parameter Group.

## Import

RDS Parameter Group can be imported using the id, e.g.

```
$ terraform import alicloud_rds_parameter_group.example <id>
```
