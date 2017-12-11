---
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_users"
sidebar_current: "docs-alicloud-datasource-ram-users"
description: |-
    Provides a list of ram users available to the user.
---

# alicloud\_ram\_users

The Ram Users data source provides a list of Alicloud Ram Users in an Alicloud account according to the specified filters.

## Example Usage

```
data "alicloud_ram_users" "user" {
  output_file = "users.txt"
  group_name = "group1"
  policy_name = "AliyunACSDefaultAccess"
  policy_type = "Custom"
  name_regex = "^user"
}

```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply to the user list returned by Alicloud.
* `group_name` - (Optional) Limit search to specific the group name. Found the users which in the specified group. 
* `policy_type` - (Optional) Limit search to specific the policy type. Valid items are `Custom` and `System`. If you set this parameter, you must set `policy_name` at one time.
* `policy_name` - (Optional) Limit search to specific the policy name. If you set this parameter without set `policy_type`, we will specified it as `System`. Found the users which attached with the specified policy.
* `output_file` - (Optional) The name of file that can save users data source after running `terraform plan`.

## Attributes Reference

A list of users will be exported and its every element contains the following attributes:

* `id` - Id of the user.
* `name` - Name of the user.
* `create_date` - Create date of the user.
* `last_login_date` - Last login date of the user.