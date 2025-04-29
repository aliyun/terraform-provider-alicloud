---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_parameter_group"
sidebar_current: "docs-alicloud-resource-polardb-parameter-group"
description: |-
  Provides a Alicloud PolarDB Parameter Group resource.
---

# alicloud_polardb_parameter_group

Provides a PolarDB Parameter Group resource.

For information about PolarDB Parameter Group and how to use it, see [What is Parameter Group](https://www.alibabacloud.com/help/en/polardb/polardb-for-mysql/user-guide/apply-a-parameter-template).

-> **NOTE:** Available since v1.183.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_parameter_group&exampleId=03f8e183-ae35-65ce-d46f-a37a51b9ab387168091d&activeTab=example&spm=docs.r.polardb_parameter_group.0.03f8e183ae&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the PolarDB Parameter Group.
* `delete` - (Defaults to 1 mins) Used when delete the PolarDB Parameter Group.

## Import

PolarDB Parameter Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_parameter_group.example <id>
```