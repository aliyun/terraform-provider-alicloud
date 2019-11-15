---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_policies"
sidebar_current: "docs-alicloud-datasource-ram-policies"
description: |-
    Provides a list of ram policies available to the user.
---

# alicloud\_ram\_policies

This data source provides a list of RAM policies in an Alibaba Cloud account according to the specified filters.

## Example Usage

```
data "alicloud_ram_policies" "policies_ds" {
  output_file = "policies.txt"
  user_name   = "user1"
  group_name  = "group1"
  type        = "System"
}

output "first_policy_name" {
  value = "${data.alicloud_ram_policies.policies_ds.policies.0.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter resulting policies by name.
* `type` - (Optional) Filter results by a specific policy type. Valid values are `Custom` and `System`.
* `user_name` - (Optional) Filter results by a specific user name. Returned policies are attached to the specified user.
* `group_name` - (Optional) Filter results by a specific group name. Returned policies are attached to the specified group.
* `role_name` - (Optional) Filter results by a specific role name. Returned policies are attached to the specified role.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of ram group names.
* `policies` - A list of policies. Each element contains the following attributes:
  * `name` - Name of the policy.
  * `type` - Type of the policy.
  * `description` - Description of the policy.
  * `default_version` - Default version of the policy.
  * `create_date` - Creation date of the policy.
  * `update_date` - Update date of the policy.
  * `attachment_count` - Attachment count of the policy.
  * `document` - Policy document of the policy.
