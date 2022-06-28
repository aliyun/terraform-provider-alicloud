---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_secret_parameter"
sidebar_current: "docs-alicloud-resource-oos-secret-parameter"
description: |-
  Provides a Alicloud OOS Secret Parameter resource.
---

# alicloud\_oos\_secret\_parameter

Provides a OOS Secret Parameter resource.

For information about OOS Secret Parameter and how to use it, see [What is Secret Parameter](https://www.alibabacloud.com/help/en/doc-detail/183418.html).

-> **NOTE:** Available in v1.147.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_kms_keys" "default" {
  status = "Enabled"
}

resource "alicloud_kms_key" "default" {
  count                  = length(data.alicloud_kms_keys.default.ids) > 0 ? 0 : 1
  description            = var.name
  status                 = "Enabled"
  pending_window_in_days = 7
}

resource "alicloud_oos_secret_parameter" "example" {
  secret_parameter_name = "example_value"
  value                 = "example_value"
  type                  = "Secret"
  key_id                = length(data.alicloud_kms_keys.default.ids) > 0 ? data.alicloud_kms_keys.default.ids.0 : concat(alicloud_kms_key.default.*.id, [""])[0]
  description           = "example_value"
  tags = {
    Created = "TF"
    For     = "OosSecretParameter"
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
}

```

## Argument Reference

The following arguments are supported:

* `constraints` - (Optional, ForceNew) The constraints of the encryption parameter. By default, this parameter is null. Valid values:
  * `AllowedValues`: The value that is allowed for the encryption parameter. It must be an array string.
  * `AllowedPattern`: The pattern that is allowed for the encryption parameter. It must be a regular expression.
  * `MinLength`: The minimum length of the encryption parameter.
  * `MaxLength`: The maximum length of the encryption parameter.
* `description` - (Optional) The description of the encryption parameter. The description must be `1` to `200` characters in length.
* `key_id` - (Optional, ForceNew) The Customer Master Key (CMK) of Key Management Service (KMS) that is used to encrypt the parameter.
* `resource_group_id` - (Optional, Computed) The ID of the Resource Group.
* `secret_parameter_name` - (Required, ForceNew) The name of the encryption parameter.  The name must be `2` to `180` characters in length, and can contain letters, digits, hyphens (-), forward slashes (/) and underscores (_). It cannot start with `ALIYUN`, `ACS`, `ALIBABA`, `ALICLOUD`, or `OOS`.
* `type` - (Optional, ForceNew) The data type of the encryption parameter. Valid values: `Secret`.
* `value` - (Required, Sensitive) The value of the encryption parameter. The value must be `1` to `4096` characters in length.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Secret Parameter. Its value is same as `secret_parameter_name`.

## Import

OOS Secret Parameter can be imported using the id, e.g.

```
$ terraform import alicloud_oos_secret_parameter.example <secret_parameter_name>
```