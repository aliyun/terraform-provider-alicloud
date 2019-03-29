---
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_policy"
sidebar_current: "docs-alicloud-resource-ram-policy"
description: |-
  Provides a RAM Policy resource.
---

# alicloud\_ram\_policy

Provides a RAM Policy resource. 

-> **NOTE:** When you want to destroy this resource forcefully(means remove all the relationships associated with it automatically and then destroy it) without set `force`  with `true` at beginning, you need add `force = true` to configuration file and run `terraform plan`, then you can delete resource forcefully.

## Example Usage

```
# Create a new RAM Policy.
resource "alicloud_ram_policy" "policy" {
  name = "test_policy"
  statement = [
    {
      effect = "Allow"
      action = [
        "oss:ListObjects",
        "oss:GetObject"
      ]
      resource = [
        "acs:oss:*:*:mybucket",
        "acs:oss:*:*:mybucket/*"
      ]
    }
  ]
  description = "this is a policy test"
  force = true
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) Name of the RAM policy. This name can have a string of 1 to 128 characters, must contain only alphanumeric characters or hyphen "-", and must not begin with a hyphen.
* `statement` - (Optional,  Type: list, Conflicts with `document`) Statements of the RAM policy document. It is required when the `document` is not specified.
     * `resource` - (Required, Type: list) List of specific objects which will be authorized. The format of each item in this list is `acs:${service}:${region}:${account_id}:${relative_id}`, such as `acs:ecs:*:*:instance/inst-002` and `acs:oss:*:1234567890000:mybucket`. The `${service}` can be `ecs`, `oss`, `ots` and so on, the `${region}` is the region info which can use `*` replace when it is not supplied, the `${account_id}` refers to someone's Alicloud account id or you can use `*` to replace, the `${relative_id}` is the resource description section which related to the `${service}`.
     * `action` - (Required, Type: list) List of operations for the `resource`. The format of each item in this list is `${service}:${action_name}`, such as `oss:ListBuckets` and `ecs:Describe*`. The `${service}` can be `ecs`, `oss`, `ots` and so on, the `${action_name}` refers to the name of an api interface which related to the `${service}`.
     * `effect` - (Required) This parameter indicates whether or not the `action` is allowed. Valid values are `Allow` and `Deny`.
* `version` - (Optional, Conflicts with `document`) Version of the RAM policy document. Valid value is `1`. Default value is `1`.
* `document` - (Optional, Conflicts with `statement` and `version`) Document of the RAM policy. It is required when the `statement` is not specified.
* `description` - (Optional, ForceNew) Description of the RAM policy. This name can have a string of 1 to 1024 characters.
* `force` - (Optional) This parameter is used for resource destroy. Default value is `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The policy ID.
* `name` - The policy name.
* `type` - The policy type.
* `description` - The policy description.
* `statement` - List of statement of the policy document.
* `document` - The policy document.
* `version` - The policy document version.
* `attachment_count` - The policy attachment count.

## Import

RAM policy can be imported using the id or name, e.g.

```
$ terraform import alicloud_ram_policy.example my-policy
```