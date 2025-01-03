---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_secret_parameter"
description: |-
  Provides a Alicloud Operation Orchestration Service (OOS) Secret Parameter resource.
---

# alicloud_oos_secret_parameter

Provides a Operation Orchestration Service (OOS) Secret Parameter resource.



For information about Operation Orchestration Service (OOS) Secret Parameter and how to use it, see [What is Secret Parameter](https://www.alibabacloud.com/help/en/doc-detail/183418.html).

-> **NOTE:** Available since v1.147.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_groups" "example" {}

resource "alicloud_kms_key" "example" {
  description            = "terraform-example"
  status                 = "Enabled"
  pending_window_in_days = 7
}

resource "alicloud_oos_secret_parameter" "example" {
  secret_parameter_name = "terraform-example"
  value                 = "terraform-example"
  type                  = "Secret"
  key_id                = alicloud_kms_key.example.id
  description           = "terraform-example"
  tags = {
    Created = "TF"
    For     = "OosSecretParameter"
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.example.groups.0.id
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
* `dkms_instance_id` - (Optional, ForceNew, Available since v1.241.0) The ID of the KMS instance.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Parameter creation time

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Secret Parameter.
* `delete` - (Defaults to 5 mins) Used when delete the Secret Parameter.
* `update` - (Defaults to 5 mins) Used when update the Secret Parameter.

## Import

Operation Orchestration Service (OOS) Secret Parameter can be imported using the id, e.g.

```shell
$ terraform import alicloud_oos_secret_parameter.example <id>
```