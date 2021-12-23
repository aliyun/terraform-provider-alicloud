---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_parameter"
sidebar_current: "docs-alicloud-resource-oos-parameter"
description: |-
  Provides a Alicloud OOS Parameter resource.
---

# alicloud\_oos\_parameter

Provides a OOS Parameter resource.

For information about OOS Parameter and how to use it, see [What is Parameter](https://www.alibabacloud.com/help/en/doc-detail/183408.html).

-> **NOTE:** Available in v1.147.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_oos_parameter" "example" {
  parameter_name = "my-Parameter"
  type           = "String"
  value          = "example_value"
  description    = "example_value"
  tags = {
    Created = "TF"
    For     = "OosParameter"
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
}

```

## Argument Reference

The following arguments are supported:

* `constraints` - (Optional, ForceNew) The constraints of the common parameter. This value follows the json format. By default, this parameter is null. Valid values:
  * `AllowedValues`: The value that is allowed for the common parameter. It must be an array string.
  * `AllowedPattern`: The pattern that is allowed for the common parameter. It must be a regular expression.
  * `MinLength`: The minimum length of the common parameter.
  * `MaxLength`: The maximum length of the common parameter.
* `description` - (Optional, Computed) The description of the common parameter. The description must be `1` to `200` characters in length.
* `parameter_name` - (Required, ForceNew) The name of the common parameter. The name must be `2` to `180` characters in length, and can contain letters, digits, hyphens (-), forward slashes (/) and underscores (_). It cannot start with `ALIYUN`, `ACS`, `ALIBABA`, `ALICLOUD`, or `OOS`.
* `resource_group_id` - (Optional, Computed) The ID of the Resource Group.
* `type` - (Required, ForceNew) The data type of the common parameter. Valid values: `String` and `StringList`.
* `value` - (Required) The value of the common parameter. The value must be `1` to `4096` characters in length.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Parameter. Its value is same as `parameter_name`.

## Import

OOS Parameter can be imported using the id, e.g.

```
$ terraform import alicloud_oos_parameter.example <parameter_name>
```