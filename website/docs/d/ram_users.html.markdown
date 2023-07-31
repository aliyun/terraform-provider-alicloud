---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_users"
sidebar_current: "docs-alicloud-datasource-ram-users"
description: |-
  Provides a list of ram users available to the user.
---

# alicloud_ram_users

This data source provides a list of RAM users in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available since v1.0.0+.

## Example Usage

```terraform
data "alicloud_ram_users" "users_ds" {
  output_file = "users.txt"
  group_name  = "group1"
  policy_name = "AliyunACSDefaultAccess"
  policy_type = "Custom"
  name_regex  = "^user"
}

output "first_user_id" {
  value = data.alicloud_ram_users.users_ds.users.0.id
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, ForceNew) A regex string to filter resulting users by their names.
* `ids` (Optional, Available since v1.53.0) - A list of ram user IDs. 
* `group_name` - (Optional, ForceNew) Filter results by a specific group name. Returned users are in the specified group. 
* `policy_type` - (Optional, ForceNew) Filter results by a specific policy type. Valid values are `Custom` and `System`. If you set this parameter, you must set `policy_name` as well.
* `policy_name` - (Optional, ForceNew) Filter results by a specific policy name. If you set this parameter without setting `policy_type`, the later will be automatically set to `System`. Returned users are attached to the specified policy.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of ram user IDs.
* `names` - A list of ram user's name. 
* `users` - A list of users. Each element contains the following attributes:
  * `id` - The original id is user's name, but it is user id in 1.37.0+.
  * `name` - Name of the user.
  * `create_date` - Creation date of the user.
  * `last_login_date` - (Removed) Last login date of the user. Removed from version 1.79.0.
