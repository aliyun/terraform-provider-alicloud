---
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_groups"
sidebar_current: "docs-alicloud-datasource-ram-groups"
description: |-
    Provides a list of ram groups available to the user.
---

# alicloud\_ram\_groups

The Ram Groups data source provides a list of Alicloud Ram Groups in an Alicloud account according to the specified filters.

## Example Usage

```
data "alicloud_ram_groups" "group" {
  output_file = "groups.txt"
  user_name = "user1"
  name_regex = "^group[0-9]*"
}

```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply to the group list returned by Alicloud.
* `user_name` - (Optional) Limit search to specific the user name. Found the groups for the specified user.
* `policy_type` - (Optional) Limit search to specific the policy type. Valid items are `Custom` and `System`. If you set this parameter, you must set `policy_name` at one time.
* `policy_name` - (Optional) Limit search to specific the policy name. If you set this parameter without set `policy_type`, we will specified it as `System`. Found the groups which attached with the specified policy.
* `output_file` - (Optional) The name of file that can save groups data source after running `terraform plan`.

## Attributes Reference

A list of groups will be exported and its every element contains the following attributes:

* `name` - Name of the group.
* `comments` - Comments of the group.