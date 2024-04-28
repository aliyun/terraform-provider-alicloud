---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_policy"
sidebar_current: "docs-alicloud-resource-ram-policy"
description: |-
  Provides a RAM Policy resource.
---

# alicloud_ram_policy

Provides a RAM Policy resource. 

For information about RAM Policy and how to use it, see [What is Policy](https://www.alibabacloud.com/help/en/ram/developer-reference/api-ram-2015-05-01-createpolicy).

-> **NOTE:** Available since vv1.0.0.

-> **NOTE:** When you want to destroy this resource forcefully(means remove all the relationships associated with it automatically and then destroy it) without set `force`  with `true` at beginning, you need add `force = true` to configuration file and run `terraform plan`, then you can delete resource forcefully.

-> **NOTE:** Each policy can own at most 5 versions and the oldest version will be removed after its version achieves 5.

-> **NOTE:** If the policy has multiple versions, all non-default versions will be deleted first when deleting policy.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ram_policy" "default" {
  policy_name     = "policyName"
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

* `policy_name` - (Optional, ForceNew, Available since v1.114.0) The name of the policy. The name must be 1 to 128 characters in length, and can contain letters, digits, and hyphens (-). It is required when the `name` is not specified. **NOTE:** If `name` is not set, `policy_name` is required.
* `policy_document` - (Optional, Available since v1.114.0) The document of the policy. see [how to use it](https://www.alibabacloud.com/help/en/ram/user-guide/policy-elements). **NOTE:** If `document` or `statement` is not set, `policy_document` is required.
* `description` - (Optional, ForceNew) The description of the policy. The description must be 1 to 1,024 characters in length.
* `rotate_strategy` - (Optional, Available since v1.114.0) The rotation strategy of the policy. The rotation strategy can be used to delete an early policy version. Default value: `None`. Valid values: `None`, `DeleteOldestNonDefaultVersionWhenLimitExceeded`.
* `force` - (Optional, Bool) Whether to forcibly delete the policy. Default value: `false`.
* `name` - (Optional, Deprecated since v1.114.0) The name of the policy. **NOTE:** Field `name` has been deprecated from provider version 1.114.0. New field `policy_name` instead.
* `document` - (Optional, Deprecated since v1.114.0) The document of the policy. see [how to use it](https://www.alibabacloud.com/help/en/ram/user-guide/policy-elements). **NOTE:** Field `document` has been deprecated from provider version 1.114.0. New field `policy_document` instead.
* `version` - (Optional, Deprecated since v1.49.0) The version of the policy document. Valid values: `1`. Default value: `1`. **NOTE:** Field `version` has been deprecated from provider version 1.49.0, and use field `document` to replace.
* `statement` - (Optional, Set, Deprecated since v1.49.0) The statement of the policy document. See [`statement`](#statement) below.
**NOTE:** Field `statement` has been deprecated from provider version 1.49.0, and use field `document` to replace.

### `statement`

The statement supports the following:

* `effect` - (Required) Specifies whether a statement result is an explicit allow or an explicit deny. Valid values: `Allow`, `Deny`.
* `action` - (Required, List) List of operations for the `resource`. The format of each item in this list is `${service}:${action_name}`, such as `oss:ListBuckets` and `ecs:Describe*`. The `${service}` can be `ecs`, `oss`, `ots` and so on, the `${action_name}` refers to the name of an api interface which related to the `${service}`.
* `resource` - (Required, List) List of specific objects which will be authorized. The format of each item in this list is `acs:${service}:${region}:${account_id}:${relative_id}`, such as `acs:ecs:*:*:instance/inst-002` and `acs:oss:*:1234567890000:mybucket`. The `${service}` can be `ecs`, `oss`, `ots` and so on, the `${region}` is the region info which can use `*` replace when it is not supplied, the `${account_id}` refers to someone's Alicloud account id or you can use `*` to replace, the `${relative_id}` is the resource description section which related to the `${service}`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Policy.
* `type` - The type of the policy.
* `version_id` - The ID of the default policy version.
* `default_version` - The default version ID of the policy.
* `attachment_count` - The number of references to the policy.

## Timeouts

-> **NOTE:** Available since v1.114.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `delete` - (Defaults to 26 mins) Used when delete the Policy.

## Import

RAM Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_ram_policy.example <id>
```
