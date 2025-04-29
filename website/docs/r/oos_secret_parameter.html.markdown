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

-> **NOTE:** Available since v1.147.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oos_secret_parameter&exampleId=8c543ccd-d749-f61b-98d8-0fa20346ebaebdb4d661&activeTab=example&spm=docs.r.oos_secret_parameter.0.8c543ccdd7&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Secret Parameter.
* `delete` - (Defaults to 5 mins) Used when delete the Secret Parameter.
* `update` - (Defaults to 5 mins) Used when update the Secret Parameter.

## Import

Operation Orchestration Service (OOS) Secret Parameter can be imported using the id, e.g.

```shell
$ terraform import alicloud_oos_secret_parameter.example <id>
```