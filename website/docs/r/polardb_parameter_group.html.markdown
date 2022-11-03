---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_parameter_group"
sidebar_current: "docs-alicloud-resource-polardb-parameter-group"
description: |-
  Provides a Alicloud PolarDB Parameter Group resource.
---

# alicloud\_polardb\_parameter\_group

Provides a PolarDB Parameter Group resource.

For information about PolarDB Parameter Group and how to use it, see [What is Parameter Group](https://www.alibabacloud.com/help/en/polardb-for-mysql/latest/createparametergroup).

-> **NOTE:** Available in v1.183.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_polardb_parameter_group" "example" {
  name       = "example_value"
  db_type    = "MySQL"
  db_version = "8.0"
  parameters {
    param_name  = "wait_timeout"
    param_value = "86400"
  }
  description = "example_value"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the parameter template. It must be 8 to 64 characters in length, and can contain letters, digits, and underscores (_). It must start with a letter and cannot contain Chinese characters.
* `db_type` - (Required, ForceNew) The type of the database engine. Only `MySQL` is supported.
* `db_version` - (Required, ForceNew) The version number of the database engine. Valid values: `5.6`, `5.7`, `8.0`.
* `parameters` - (Required, ForceNew) The parameter template. See the following `Block parameters`.
* `description` - (Optional, ForceNew) The description of the parameter template. It must be 0 to 200 characters in length.

#### Block parameters

The parameters supports the following:

* `param_name` - (Required, ForceNew) The name of a parameter in the parameter template.
* `param_value` - (Required, ForceNew) The value of a parameter in the parameter template.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Parameter Group.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the PolarDB Parameter Group.
* `delete` - (Defaults to 1 mins) Used when delete the PolarDB Parameter Group.

## Import

PolarDB Parameter Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_parameter_group.example <id>
```