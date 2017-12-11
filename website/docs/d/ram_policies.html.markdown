---
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_policies"
sidebar_current: "docs-alicloud-datasource-ram-policies"
description: |-
    Provides a list of ram policies available to the user.
---

# alicloud\_ram\_policies

The Ram Policies data source provides a list of Alicloud Ram Policies in an Alicloud account according to the specified filters.

## Example Usage

```
data "alicloud_ram_policies" "policy" {
  output_file = "policies.txt"
  user_name = "user1"
  group_name = "group1"
  type = "System"
}

```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply to the policy list returned by Alicloud.
* `type` - (Optional) Limit search to specific the policy type. Valid items are `Custom` and `System`.
* `user_name` - (Optional) Limit search to specific the user name. Found the policies which attached with the specified user.
* `group_name` - (Optional) Limit search to specific the group name. Found the policies which attached with the specified group.
* `role_name` - (Optional) Limit search to specific the role name. Found the policies which attached with the specified role.
* `output_file` - (Optional) The name of file that can save policies data source after running `terraform plan`.

## Attributes Reference

A list of policies will be exported and its every element contains the following attributes:

* `name` - Name of the policy.
* `type` - Type of the policy.
* `description` - Description of the policy.
* `default_version` - Default version of the policy.
* `create_date` - Create date of the policy.
* `update_date` - Update date of the policy.
* `attachment_count` - Attachment count of the policy.
* `document` - Policy document of the policy.