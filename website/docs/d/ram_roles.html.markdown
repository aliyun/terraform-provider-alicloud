---
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_roles"
sidebar_current: "docs-alicloud-datasource-ram-roles"
description: |-
    Provides a list of ram roles available to the user.
---

# alicloud\_ram\_roles

The Ram Roles data source provides a list of Alicloud Ram Roles in an Alicloud account according to the specified filters.

## Example Usage

```
data "alicloud_ram_roles" "role" {
  output_file = "roles.txt"
  name_regex = ".*test.*"
  policy_name = "AliyunACSDefaultAccess"
  policy_type = "Custom"
}

```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply to the role list returned by Alicloud.
* `policy_type` - (Optional) Limit search to specific the policy type. Valid items are `Custom` and `System`. If you set this parameter, you must set `policy_name` at one time.
* `policy_name` - (Optional) Limit search to specific the policy name. If you set this parameter without set `policy_type`, we will specified it as `System`. Found the roles which attached with the specified policy.
* `output_file` - (Optional) The name of file that can save roles data source after running `terraform plan`.

## Attributes Reference

A list of roles will be exported and its every element contains the following attributes:

* `id` - Id of the role.
* `name` - Name of the role.
* `arn` - Resource descriptor of the role.
* `description` - Description of the role.
* `assume_role_policy_document` - Authorization strategy of the role. This parameter is deprecated and replaced by `document`.
* `document` - Authorization strategy of the role.
* `create_date` - Create date of the role.
* `update_date` - Update date of the role.