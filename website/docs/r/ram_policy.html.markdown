---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_policy"
description: |-
  Provides a Alicloud RAM Policy resource.
---

# alicloud_ram_policy

Provides a RAM Policy resource.



For information about RAM Policy and how to use it, see [What is Policy](https://www.alibabacloud.com/help/en/ram/developer-reference/api-ram-2015-05-01-createpolicy).

-> **NOTE:** Available since v1.0.0.

-> **NOTE:** When you want to destroy this resource forcefully(means remove all the relationships associated with it automatically and then destroy it) without set `force`  with `true` at beginning, you need add `force = true` to configuration file and run `terraform plan`, then you can delete resource forcefully.

-> **NOTE:** Each policy can own at most 5 versions and the oldest version will be removed after its version achieves 5.

-> **NOTE:** If the policy has multiple versions, all non-default versions will be deleted first when deleting policy.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ram_policy&exampleId=8efe2170-27c1-b4fc-82d2-b2fff764cc1d424c720a&activeTab=example&spm=docs.r.ram_policy.0.8efe217027&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
# Create a new RAM Policy.
resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_ram_policy" "policy" {
  policy_name     = "tf-example-${random_integer.default.result}"
  policy_document = <<EOF
  {
    "Statement": [
      {
        "Action": [
          "oss:ListObjects",
          "oss:GetObject"
        ],
        "Effect": "Allow",
        "Resource": [
          "acs:oss:*:*:mybucket",
          "acs:oss:*:*:mybucket/*"
        ]
      }
    ],
      "Version": "1"
  }
  EOF
  description     = "this is a policy test"
}
```
## Argument Reference

The following arguments are supported:
* `description` - (Optional) The description of the policy. The description must be `1` to `1024` characters in length. **NOTE:** From version 1.246.0, `description` can be modified.
* `policy_document` - (Optional, Available since v1.114.0) The document of the policy. It is required when the `statement` is not specified.
* `policy_name` - (Optional, ForceNew) The name of the policy. The name must be `1` to `128` characters in length and can contain letters, digits, and hyphens (-). It is required when the `name` is not specified.
* `rotate_strategy` - (Optional, Available since v1.114.0) The rotation strategy of the policy. The rotation strategy can be used to delete an early policy version. Default to `None`. Valid values: `None`, `DeleteOldestNonDefaultVersionWhenLimitExceeded`.
* `force` - (Optional, Bool) Specifies whether to force delete the Policy. Default value: `false`. Valid values:
  - `true`: Enable.
  - `false`: Disable.
* `tags` - (Optional, Map, Available since v1.246.0) A mapping of tags to assign to the resource.
-> **NOTE:** Currently, `tags` is available on `Domestic Site Account`.
* `name` - (Optional, ForceNew, Deprecated since v1.114.0) Field `name` has been deprecated from provider version 1.114.0. New field `policy_name` instead.
* `document` - (Optional, Deprecated since v1.114.0) Field `document` has been deprecated from provider version 1.114.0. New field `policy_document` instead.
* `version` - (Optional, Deprecated since v1.49.0) Field `version` has been deprecated from provider version 1.49.0. New field `document` instead.
* `statement` - (Optional, List, Deprecated since v1.49.0) Field `statement` has been deprecated from provider version 1.49.0. New field `document` instead. See [`statement`](#statement) below.

### `statement`

The statement support the following:
* `resource` - (Deprecated since 1.49.0, Required, Type: list) (It has been deprecated since version 1.49.0, and use field `document` to replace.) List of specific objects which will be authorized. The format of each item in this list is `acs:${service}:${region}:${account_id}:${relative_id}`, such as `acs:ecs:*:*:instance/inst-002` and `acs:oss:*:1234567890000:mybucket`. The `${service}` can be `ecs`, `oss`, `ots` and so on, the `${region}` is the region info which can use `*` replace when it is not supplied, the `${account_id}` refers to someone`s Alicloud account id or you can use `*` to replace, the `${relative_id}` is the resource description section which related to the `${service}`.
* `action` - (Deprecated since 1.49.0, Required, Type: list) (It has been deprecated since version 1.49.0, and use field `document` to replace.) List of operations for the `resource`. The format of each item in this list is `${service}:${action_name}`, such as `oss:ListBuckets` and `ecs:Describe*`. The `${service}` can be `ecs`, `oss`, `ots` and so on, the `${action_name}` refers to the name of an api interface which related to the `${service}`.
* `effect` - (Deprecated since 1.49.0, Required) (It has been deprecated since version 1.49.0, and use field `document` to replace.) This parameter indicates whether or not the `action` is allowed. Valid values are `Allow` and `Deny`.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Policy.
* `type` - The type of the policy.
* `attachment_count` - The number of references to the policy.
* `version_id` - The ID of the default policy version.
* `default_version` - The default version ID of the policy.
* `create_time` - (Available since v1.246.0) The time when the policy was created.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Policy.
* `update` - (Defaults to 5 mins) Used when update the Policy.

## Import

RAM Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_ram_policy.example <id>
```
