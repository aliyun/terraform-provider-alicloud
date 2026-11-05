---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_parameter"
description: |-
  Provides a Alicloud Data Works Parameter resource.
---

# alicloud_data_works_parameter

Provides a Data Works Parameter resource.

DataWorks parameters are managed constants and variables reused across nodes in a workspace. A parameter with `type = SecretConstant` is encrypted at rest; the production environment (`Prod`) configuration is required for it.

For information about Data Works Parameter and how to use it, see [What is Parameter](https://www.alibabacloud.com/help/en/dataworks/developer-reference/api-dataworks-public-2024-05-18-createparameter).

-> **NOTE:** Available since v1.245.0.

-> **NOTE:** The DataWorks Parameter API requires DataWorks Professional Edition or above.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}

resource "random_integer" "default" {
  max = 99999
  min = 10000
}

data "alicloud_account" "default" {}

resource "alicloud_data_works_project" "default" {
  project_name            = "${var.name}${random_integer.default.result}"
  display_name           = var.name
  description            = "tf_for_parameter"
  dev_environment_enabled = true
}

resource "alicloud_data_works_parameter" "default" {
  project_id = alicloud_data_works_project.default.id
  name       = "workspace.${var.name}${random_integer.default.result}"
  type       = "PlainConstant"
  owner      = data.alicloud_account.default.id
  description = "managed by terraform"

  properties {
    env_type = "Prod"
    value    = "prod_value"
  }

  properties {
    env_type = "Dev"
    value    = "dev_value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The parameter name, unique within the workspace. It must start with `workspace.`. The remaining part must start with a letter and can only contain letters, underscores and digits. The total length cannot exceed 255 characters.
* `type` - (Required, ForceNew) The type of the parameter. Valid values: `PlainConstant` (constant, plaintext), `SecretConstant` (constant, ciphertext) and `Variable` (variable).
* `project_id` - (Required, ForceNew) The workspace ID. Required when `scope = Project`.
* `owner` - (Required) The account ID of the responsible person.
* `properties` - (Required) The value configuration of the parameter. The production environment (`Prod`) is required. If duplicate environment configurations appear in the array, the later configuration does not take effect. See [`properties`](#properties) below.
* `description` - (Optional) The description of the parameter.
* `scope` - (Optional, ForceNew) The scope of the parameter. Valid value: `Project`. Defaults to `Project`. Currently no other type can be selected.

### properties

The `properties` block supports:

* `env_type` - (Required) The project environment. Valid values: `Prod` (production) and `Dev` (development). A `Prod` configuration is required for `type = SecretConstant`.
* `value` - (Required) The parameter value. Only Chinese characters, letters, digits and some special characters (`/`, `:`, `.`, `[`, `]`, `,`, `\`, `'`, `"`, `_`, `=`, `?`, space, carriage return, line feed, `+`, `-`, `*`, `%`, `&`, `@`, `!`, `$`, `#`, `{`, `}`) are allowed.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the DataWorks Parameter.
* `create_time` - The creation time. A 13-digit millisecond timestamp, for example `1640000000000`.
* `create_user` - The creator account ID.
* `modify_time` - The last modification time. A 13-digit millisecond timestamp, for example `1640000000000`.
* `modify_user` - The last modifier account ID.
* `version` - The version number.

## Import

Data Works Parameter can be imported using the id, e.g.

```shell
$ terraform import alicloud_data_works_parameter.example <parameter_id>
```
